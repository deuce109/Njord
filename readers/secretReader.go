package readers

import (
	"os"
	"path/filepath"
)

func ReadSecretFile(secretFilepath string) (string, error) {
	cleanPath := filepath.Clean(secretFilepath)
	dataBytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", err
	} else {
		return string(dataBytes), nil
	}
}
