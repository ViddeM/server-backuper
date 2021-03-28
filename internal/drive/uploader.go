package drive

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

var folderMimeType = "application/vnd.google-apps.folder"

func UploadFile(driveId, filePath string, currentTime *time.Time) {
	ctx := context.Background()
	service, err := drive.NewService(ctx, option.WithCredentialsFile("service-account.json"),
		option.WithScopes("https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.readonly"))

	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	aboutService := drive.NewAboutService(service)
	about, err := aboutService.Get().Fields("user/displayName").Do()
	if err != nil {
		log.Fatalf("Failed to retrieve about: %v", err)
	}
	log.Printf("Accessing drive as user %s\n", about.User.DisplayName)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Failed to open file test.txt")
	}

	fileService := drive.NewFilesService(service)
	files, err := fileService.List().
		SupportsAllDrives(true).
		Corpora("drive").
		IncludeItemsFromAllDrives(true).
		Q(fmt.Sprintf("mimeType='%s'", folderMimeType)).
		DriveId(driveId).
		Fields("files(parents,name,id,trashed)").
		Do()

	if err != nil {
		log.Fatalf("Failed to retrieve files list: %v\n", err)
	}

	parentFolder := getOrCreateParentFolder(files.Files, currentTime, driveId, fileService)

	uploadedFile, err := fileService.Create(&drive.File{
		DriveId:  driveId,
		Name:     filePath,
		MimeType: "application/zip",
		Parents:  []string{parentFolder.Id},
	}).SupportsAllDrives(true).
		Media(f).
		Do()
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	log.Printf("Uploaded file '%s' to '%s'\n", uploadedFile.Name, uploadedFile.DriveId)
}

func getOrCreateParentFolder(files []*drive.File, currentTime *time.Time, driveId string, service *drive.FilesService) *drive.File {
	backups := findChildWithName(driveId, "backups", files)
	if backups == nil {
		log.Fatalf("Failed to find backups folder")
	}

	yearName := fmt.Sprintf("%04d", currentTime.Year())
	monthName := fmt.Sprintf("%02d", currentTime.Month())
	dayName := fmt.Sprintf("%02d", currentTime.Day())
	year := findChildWithName(backups.Id, yearName, files)
	if year == nil {
		year = createFolder(yearName, backups.Id, driveId, service)
		month := createFolder(monthName, year.Id, driveId, service)
		day := createFolder(dayName, month.Id, driveId, service)
		return day
	}

	month := findChildWithName(year.Id, monthName, files)
	if month == nil {
		month := createFolder(monthName, year.Id, driveId, service)
		day := createFolder(dayName, month.Id, driveId, service)
		return day
	}

	day := findChildWithName(month.Id, dayName, files)
	if day == nil {
		day := createFolder(dayName, month.Id, driveId, service)
		return day
	}

	return day
}

func createFolder(name, parentId, driveId string, service *drive.FilesService) *drive.File {
	file, err := service.Create(&drive.File{
		DriveId: driveId,
		Name: name,
		MimeType: folderMimeType,
		Parents: []string{parentId},
	}).SupportsAllDrives(true).
		Do()

	if err != nil {
		log.Fatalf("Failed to create directory %s: %v\n", name, err)
	}

	return file
}

func findChildWithName(expectedParent, folderName string, files []*drive.File) *drive.File {
	for _, file := range files {
		if !file.Trashed && file.Name == folderName {
			for _, parent := range file.Parents {
				if parent == expectedParent {
					return file
				}
			}
		}
	}
	return nil
}