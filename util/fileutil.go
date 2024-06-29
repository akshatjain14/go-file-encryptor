package util

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const key = "ABCDEFGHIJKLMNOPQRSTUVWXYS123456"

func ScanDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func SaveHashMap(filePath string, hashMap map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for path, hash := range hashMap {
		_, err := writer.WriteString(fmt.Sprintf("%s %s\n", path, hash))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func LoadHashMap(filePath string) (map[string]string, error) {
	hashMap := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) == 2 {
			hashMap[parts[0]] = parts[1]
		}
	}
	return hashMap, scanner.Err()
}

func CheckAndDecrypt(filePath, hashFilePath, destPath string) error {
	hashMap, err := LoadHashMap(hashFilePath)
	if err != nil {
		return err
	}

	expectedHash, ok := hashMap[filePath]
	if !ok {
		return fmt.Errorf("hash not found for file %s", filePath)
	}

	hashValue, err := computeHash(filePath)
	if err != nil {
		return err
	}

	if hashValue != expectedHash {
		return fmt.Errorf("hash mismatch for file %s", filePath)
	}

	return decryptFile(filePath, destPath)
}

func decryptFile(filePath, destPath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	if len(content) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short")
	}

	iv := content[:aes.BlockSize]
	content = content[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(content, content)

	return os.WriteFile(destPath, content, 0644)
}

func computeHash(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), nil
}
