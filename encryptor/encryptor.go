package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
)

const key = "ABCDEFGHIJKLMNOPQRSTUVWXYS123456"

func EncryptFile(sourcePath, destDir string) (string, error) {
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(content))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], content)

	destPath := filepath.Join(destDir, filepath.Base(sourcePath))
	err = os.WriteFile(destPath, ciphertext, 0644)
	if err != nil {
		return "", err
	}

	return destPath, nil
}
