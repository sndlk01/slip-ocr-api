package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
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

// CleanOCRText cleans and formats OCR text for better readability
func CleanOCRText(text string) string {
	// Remove excessive whitespace and clean up the text
	lines := strings.Split(text, "\n")
	var cleanedLines []string

	for _, line := range lines {
		// Trim leading and trailing whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Replace multiple spaces with single space
		re := regexp.MustCompile(`\s+`)
		line = re.ReplaceAllString(line, " ")

		cleanedLines = append(cleanedLines, line)
	}

	// Join lines with newline
	return strings.Join(cleanedLines, "\n")
}
