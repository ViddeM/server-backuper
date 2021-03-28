package crypto

import (
	"io/ioutil"
	"log"
)

func EncryptFile(toEncrypt string, key []byte) string {
	k := getKey(key)

	data, err := ioutil.ReadFile(toEncrypt)
	if err != nil {
		log.Fatalf("Failed to read file to encrypt: %v", err)
	}

	encrypted, err := k.Encrypt(data)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}

	outputFile := getEncryptedFilename(toEncrypt)
	err = ioutil.WriteFile(outputFile, encrypted, 0644)
	if err != nil {
		log.Fatalf("Failed to write encrypted data to file: %v", err)
	}

	return outputFile
}
