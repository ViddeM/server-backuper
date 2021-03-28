package main

import (
	"github.com/viddem/server-backuper/internal/arguments"
	"github.com/viddem/server-backuper/internal/drive"
	"github.com/viddem/server-backuper/internal/zip"
	"time"

	"log"
	"os"
)

func main() {
	args := arguments.LoadArgs()
	currentTime := time.Now()
	file, err := zip.Directory(args.BackupFolder, args.LocalStorageFolder, &currentTime)
	if err != nil {
		log.Fatalf("Failed to create zip: %v", err)
	}
	log.Printf("Successfully created zip: '%s'", file)

	//outputFile := crypto.EncryptFile(file, args.Key)
	//log.Printf("Encrypted file to '%s'\n", outputFile)

	drive.UploadFile(args.BackupDriveId, file, &currentTime)
	deleteLocalBackups(args.LocalStorageFolder)
}

func deleteLocalBackups(backupPath string) {
	err := os.RemoveAll(backupPath)
	if err != nil {
		log.Fatalf("Failed to clear local backups folder: %v", err)
	}

	log.Println("Successfully cleared local backup folder!")
}
