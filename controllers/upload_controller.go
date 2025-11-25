package controllers

import (
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

	file, err := ctx.FormFile("slip")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded. Use 'slip' as the form field name",
		})
		return
	}

	if file.Size > config.AppConfig.MaxUploadSize {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "File too large. Maximum size is 10MB",
		})
		return
	}

	filename := strings.ToLower(file.Filename)
	if err := c.ocrService.ValidateImage(filename); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	uniqueFilename := utils.GenerateUniqueFilename(filename)
	uploadPath := filepath.Join(config.AppConfig.UploadDir, uniqueFilename)

	if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
		log.Printf("Failed to save uploaded file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save uploaded file",
		})
		return
	}

	log.Printf("File uploaded: %s (%.2f KB)", uniqueFilename, float64(file.Size)/1024)

	transaction, err := c.ocrService.ProcessSlip(uploadPath, req.Type)
	if err != nil {
		c.ocrService.CleanupUploadedFile(uploadPath)

		log.Printf("OCR processing failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process slip: " + err.Error(),
		})
		return
	}

	if err := c.transactionService.Create(transaction); err != nil {
		c.ocrService.CleanupUploadedFile(uploadPath)

		log.Printf("Failed to save transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save transaction",
		})
		return
	}

	if err := c.ocrService.CleanupUploadedFile(uploadPath); err != nil {
		log.Printf("Warning: failed to cleanup uploaded file: %v", err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Slip processed successfully",
		"transaction": transaction,
	})
}
