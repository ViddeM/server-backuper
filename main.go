package main

import (
	"archive/zip"
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProgramArguments struct {
	BackupDriveId string
	BackupFolder string
	LocalStorageFolder string
}

func loadArgs() *ProgramArguments {
	args := os.Args[1:]
	driveId := args[0]
	backupFolder := args[1]

	fileInfo, err := os.Stat(backupFolder)
	if err != nil {
		log.Fatalf("Failed retrieving info about backupFodler '%s', err: '%v'\n", backupFolder, err)
	}

	if fileInfo.IsDir() == false {
		log.Fatalf("Backup folder must be a directory, '%s' is not\n", backupFolder)
	}

	localStorageFolder := "./backups"

	return &ProgramArguments{
		BackupDriveId: driveId,
		BackupFolder: backupFolder,
		LocalStorageFolder: localStorageFolder,
	}
}

func main() {
	args := loadArgs()
	file, err := zipFolder(args.BackupFolder, args.LocalStorageFolder)
	if err != nil {
		log.Fatalf("Failed to create zip: %v", err)
	}
	log.Printf("Should have created zipfile: '%s'", file)
	uploadFile(args.BackupDriveId, file)
}

func zipFolder(path, outputPath string) (string, error) {
	currentTime := time.Now()

	CreatePathIfNotExist(outputPath)

	zipName := fmt.Sprintf("%s/backup_%s.zip", outputPath, currentTime.Format("2006-01-02_15:04:05.000"))
	zipDir, err := os.Create(zipName)
	if err != nil {
		log.Fatalf("Failed to create zip: %v", err)
	}

	zipWriter := zip.NewWriter(zipDir)
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(path))
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	err = zipWriter.Close()
	return zipDir.Name(), nil
}

func CreatePathIfNotExist(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				log.Fatalf("Failed to create directory: %v", err)
			}
			return
		} else {
			log.Fatalf("Failed to stat path '%s', err: %v", path, err)
		}
 	} else if fileInfo.IsDir() == false {
		log.Fatalf("Local backup path ('%s') must be a directory!", path)
	}
}

func uploadFile(driveId, filePath string) {
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
	uploadedFile, err := fileService.Create(&drive.File{
		DriveId:  driveId,
		Name:     filePath,
		MimeType: "application/zip",
		Parents:  []string{driveId},
	}).SupportsAllDrives(true).SupportsTeamDrives(true).Media(f).Do()
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	log.Printf("Uploaded file '%s' to '%s'\n", uploadedFile.Name, uploadedFile.DriveId)
}
