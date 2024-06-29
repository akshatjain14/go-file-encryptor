package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"go-file-encryptor/encryptor"
	"go-file-encryptor/hash"
	"go-file-encryptor/util"
)

const (
	sourceDir      = "data/source/"
	destinationDir = "data/destination/"
	hashFilePath   = "data/hashes.txt"
)

func main() {
	// Convert to absolute paths
	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		log.Fatalf("Error converting source directory to absolute path: %v", err)
	}

	absDestinationDir, err := filepath.Abs(destinationDir)
	if err != nil {
		log.Fatalf("Error converting destination directory to absolute path: %v", err)
	}

	absHashFilePath, err := filepath.Abs(hashFilePath)
	if err != nil {
		log.Fatalf("Error converting hash file path to absolute path: %v", err)
	}

	files, err := util.ScanDirectory(absSourceDir)
	if err != nil {
		log.Fatalf("Error scanning source directory: %v", err)
	}

	hashMap := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go processFile(file, absDestinationDir, &wg, &mu, hashMap)
	}

	wg.Wait()

	err = util.SaveHashMap(absHashFilePath, hashMap)
	if err != nil {
		log.Fatalf("Error saving hash map: %v", err)
	}

	// Example integrity check and decryption
	testFile := "example.txt"
	encryptedFile := filepath.Join(absDestinationDir, testFile)
	err = util.CheckAndDecrypt(encryptedFile, absHashFilePath, filepath.Join("data/source", testFile))
	if err != nil {
		log.Fatalf("Error checking and decrypting file: %v", err)
	}

	fmt.Println("Process completed successfully!")
}

func processFile(file, absDestinationDir string, wg *sync.WaitGroup, mu *sync.Mutex, hashMap map[string]string) {
	defer wg.Done()

	encryptedFilePath, err := encryptor.EncryptFile(file, absDestinationDir)
	if err != nil {
		log.Printf("Error encrypting file %s: %v", file, err)
		return
	}

	hashValue, err := hash.ComputeHash(encryptedFilePath)
	if err != nil {
		log.Printf("Error computing hash for file %s: %v", encryptedFilePath, err)
		return
	}

	mu.Lock()
	hashMap[encryptedFilePath] = hashValue
	mu.Unlock()
}
