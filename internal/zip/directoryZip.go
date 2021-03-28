package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Directory(path, outputPath string, currentTime *time.Time) (string, error) {

	createPathIfNotExist(outputPath)

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

func createPathIfNotExist(path string) {
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

