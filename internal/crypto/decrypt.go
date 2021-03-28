package crypto

import (
	"io/ioutil"
	"log"
)

func DecryptFile(toDecrypt string, key []byte) string {
	k := getKey(key)

	data, err := ioutil.ReadFile(toDecrypt)
	if err != nil {
		log.Fatalf("Failed to read file to decrypt: %v", err)
	}

	decrypted, err := k.Decrypt(data)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}

	outputFile := getDecryptedFilename(toDecrypt)
	err = ioutil.WriteFile(outputFile, decrypted, 0644)
	if err != nil {
		log.Fatalf("Failed to write decrypted data to file: %v", err)
	}

	return outputFile
}