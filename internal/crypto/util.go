package crypto

import (
	"fmt"
	"github.com/ovh/symmecrypt"
	"github.com/ovh/symmecrypt/keyloader"
	"log"
	"path/filepath"
	"time"
)

func getKey(key []byte) symmecrypt.Key {
	k, err := keyloader.NewKey(&keyloader.KeyConfig{
		Identifier: "storage",
		Cipher:     "aes-gcm",
		Timestamp:  time.Now().Unix(),
		Sealed:     false,
		Key:        string(key),
	})

	if err != nil {
		log.Fatalf("Failed to create new key config: %v", err)
	}

	return k
}


func getEncryptedFilename(original string) string {
	extension := filepath.Ext(original)
	nameOnly := original[:len(original) - len(extension)]
	return fmt.Sprintf("%s%s", nameOnly, ".encrypted")
}

func getDecryptedFilename(encrypted string) string {
	extension := filepath.Ext(encrypted)
	nameOnly := encrypted[:len(encrypted) - len(extension)]
	return fmt.Sprintf("%s.zip", nameOnly)
}