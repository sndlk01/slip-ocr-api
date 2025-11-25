package ocr

import (
	"fmt"
	"log"

	"github.com/otiai10/gosseract/v2"
)

func PerformOCR(imagePath string, lang string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetLanguage(lang)
	if err != nil {
		return "", fmt.Errorf("failed to set language: %w", err)
	}

	err = client.SetImage(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to set image: %w", err)
	}

	client.SetPageSegMode(gosseract.PSM_AUTO)

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to perform OCR: %w", err)
	}

	log.Printf("OCR completed. Extracted %d characters", len(text))

	return text, nil
}
