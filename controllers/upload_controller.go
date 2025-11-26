package controllers

import (
	"fmt"
	"log"
	"net/http"
	"ocr-api/config"
	"ocr-api/services"
	"ocr-api/utils"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	ocrService         *services.OCRService
	transactionService *services.TransactionService
}

func NewUploadController() *UploadController {
	return &UploadController{
		ocrService:         services.NewOCRService(),
		transactionService: services.NewTransactionService(),
	}
}

type UploadRequest struct {
	Type string `form:"type" binding:"required"` 
}

func (c *UploadController) UploadSlip(ctx *gin.Context) {
	var req UploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing or invalid 'type' field. Must be 'income' or 'expense'",
		})
		return
	}

	if !utils.ValidateTransactionType(req.Type) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction type. Must be 'income' or 'expense'",
		})
		return
	}

	// Get multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form",
		})
		return
	}

	// Support both 'slip' (single) and 'slips' (multiple)
	files := form.File["slips"]
	if len(files) == 0 {
		files = form.File["slip"]
	}

	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No files uploaded. Use 'slip' or 'slips' as the form field name",
		})
		return
	}

	var transactions []interface{}
	var errors []string
	var uploadPaths []string

	// Process each file
	for _, file := range files {
		if file.Size > config.AppConfig.MaxUploadSize {
			errors = append(errors, fmt.Sprintf("File '%s' is too large (max 10MB)", file.Filename))
			continue
		}

		filename := strings.ToLower(file.Filename)
		if err := c.ocrService.ValidateImage(filename); err != nil {
			errors = append(errors, fmt.Sprintf("File '%s': %s", file.Filename, err.Error()))
			continue
		}

		uniqueFilename := utils.GenerateUniqueFilename(filename)
		uploadPath := filepath.Join(config.AppConfig.UploadDir, uniqueFilename)

		if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
			log.Printf("Failed to save uploaded file '%s': %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to save file '%s'", file.Filename))
			continue
		}

		uploadPaths = append(uploadPaths, uploadPath)
		log.Printf("File uploaded: %s (%.2f KB)", uniqueFilename, float64(file.Size)/1024)

		transaction, err := c.ocrService.ProcessSlip(uploadPath, req.Type)
		if err != nil {
			log.Printf("OCR processing failed for '%s': %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to process '%s': %s", file.Filename, err.Error()))
			continue
		}

		if err := c.transactionService.Create(transaction); err != nil {
			log.Printf("Failed to save transaction for '%s': %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to save transaction for '%s'", file.Filename))
			continue
		}

		transactions = append(transactions, transaction)
	}

	// Cleanup all uploaded files
	for _, uploadPath := range uploadPaths {
		if err := c.ocrService.CleanupUploadedFile(uploadPath); err != nil {
			log.Printf("Warning: failed to cleanup uploaded file: %v", err)
		}
	}

	// Return response
	if len(transactions) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  "No slips were processed successfully",
			"errors": errors,
		})
		return
	}

	response := gin.H{
		"message":      fmt.Sprintf("Processed %d out of %d slips successfully", len(transactions), len(files)),
		"transactions": transactions,
		"success_count": len(transactions),
		"total_count":   len(files),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	ctx.JSON(http.StatusCreated, response)
}
