package arguments

import (
	"io/ioutil"
	"log"
	"os"
)

type ProgramArguments struct {
	BackupDriveId string
	BackupFolder string
	LocalStorageFolder string
	Key []byte
}

func LoadArgs() *ProgramArguments {
	args := os.Args[1:]
	driveId := args[0]
	backupFolder := args[1]

	fileInfo, err := os.Stat(backupFolder)
	if err != nil {
		log.Fatalf("Failed retrieving info about backupFolder '%s', err: '%v'\n", backupFolder, err)
	}

	if fileInfo.IsDir() == false {
		log.Fatalf("Backup folder must be a directory, '%s' is not\n", backupFolder)
	}

	localStorageFolder := "./backups"

	keyFile := "./key.sec"
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Failed to read keyfile '%s', err: %v", keyFile, err)
	}

	return &ProgramArguments{
		BackupDriveId: driveId,
		BackupFolder: backupFolder,
		LocalStorageFolder: localStorageFolder,
		Key: key,
	}
}

