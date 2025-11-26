package services

import (
	"fmt"
	"log"
	"ocr-api/config"
	"ocr-api/models"
	"ocr-api/ocr"
	"ocr-api/utils"
	"os"
	"path/filepath"
)

type OCRService struct{}

func NewOCRService() *OCRService {
	return &OCRService{}
}

func (s *OCRService) ProcessSlip(imagePath string, transactionType string) (*models.Transaction, *models.Subscription, error) {
	log.Printf("Processing slip: %s", imagePath)

	jpegPath, err := ocr.ConvertToJPEG(imagePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert image: %w", err)
	}
	defer s.cleanupFile(jpegPath)

	processedPath, err := ocr.PreprocessImage(jpegPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to preprocess image: %w", err)
	}
	defer s.cleanupFile(processedPath)

	ocrText, err := ocr.PerformOCR(processedPath, config.AppConfig.TesseractLang)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to perform OCR: %w", err)
	}

	log.Printf("OCR Text:\n%s\n", ocrText)

	// Clean OCR text for better readability
	cleanedOCRText := utils.CleanOCRText(ocrText)

	extractedData, err := ocr.ExtractData(ocrText)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract data: %w", err)
	}

	normalizedDate := ocr.NormalizeDate(extractedData.Date)

	transaction := &models.Transaction{
		Type:       transactionType,
		Amount:     extractedData.Amount,
		Date:       normalizedDate,
		Time:       extractedData.Time,
		Reference:  extractedData.Reference,
		Bank:       extractedData.Bank,
		Sender:     extractedData.Sender,
		Receiver:   extractedData.Receiver,
		RawOCRText: cleanedOCRText,
	}

	// Auto-detect subscription
	subscriptionService := NewSubscriptionService()
	detectedSub := subscriptionService.DetectSubscription(ocrText, extractedData.Amount)

	log.Printf("Transaction created: %+v", transaction)

	return transaction, detectedSub, nil
}

func (s *OCRService) cleanupFile(filePath string) {
	if filePath == "" {
		return
	}

	if err := os.Remove(filePath); err != nil {
		log.Printf("Warning: failed to cleanup file %s: %v", filePath, err)
	} else {
		log.Printf("Cleaned up file: %s", filePath)
	}
}

func (s *OCRService) SaveUploadedFile(fileHeader interface{}, filename string) (string, error) {
	uploadPath := filepath.Join(config.AppConfig.UploadDir, filename)

	return uploadPath, nil
}

func (s *OCRService) CleanupUploadedFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	log.Printf("Removed uploaded file: %s", filePath)
	return nil
}

func (s *OCRService) ValidateImage(filename string) error {
	ext := filepath.Ext(filename)
	validExts := []string{".jpg", ".jpeg", ".png"}

	if !utils.Contains(validExts, ext) {
		return fmt.Errorf("invalid file type: %s. Allowed types: jpg, jpeg, png", ext)
	}

	return nil
}
