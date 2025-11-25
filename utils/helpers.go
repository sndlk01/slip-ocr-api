package utils

import (
	"fmt"
	"path/filepath"
	"time"
)

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, ext)
}

func ValidateTransactionType(transactionType string) bool {
	validTypes := []string{"income", "expense"}
	return Contains(validTypes, transactionType)
}
