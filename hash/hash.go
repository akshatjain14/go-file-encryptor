package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func ComputeHash(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), nil
}
