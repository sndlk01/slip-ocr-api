package controllers

import (
	"net/http"
	"ocr-api/models"
	"ocr-api/services"
	"ocr-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service *services.TransactionService
}

func NewTransactionController() *TransactionController {
	return &TransactionController{
		service: services.NewTransactionService(),
	}
}

func (c *TransactionController) GetAll(ctx *gin.Context) {
	transactionType := ctx.Query("type")
	bank := ctx.Query("bank")

	var transactions interface{}
	var err error

	if transactionType != "" {
		transactions, err = c.service.GetByType(transactionType)
	} else if bank != "" {
		transactions, err = c.service.GetByBank(bank)
	} else {
		transactions, err = c.service.GetAll()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}

func (c *TransactionController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	transaction, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
	})
}

func (c *TransactionController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	err = c.service.Delete(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})
}

type CreateTransactionRequest struct {
	Type      string  `json:"type" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	Reference string  `json:"reference"`
	Bank      string  `json:"bank"`
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Detail    string  `json:"detail"`
}

func (c *TransactionController) Create(ctx *gin.Context) {
	var req CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if !utils.ValidateTransactionType(req.Type) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction type. Must be 'income' or 'expense'",
		})
		return
	}

	transaction := &models.Transaction{
		Type:      req.Type,
		Amount:    req.Amount,
		Date:      req.Date,
		Time:      req.Time,
		Reference: req.Reference,
		Bank:      req.Bank,
		Sender:    req.Sender,
		Receiver:  req.Receiver,
		Detail:    req.Detail,
	}

	if err := c.service.Create(transaction); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create transaction",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

type UpdateTransactionRequest struct {
	Amount    *float64 `json:"amount"`
	Date      *string  `json:"date"`
	Time      *string  `json:"time"`
	Reference *string  `json:"reference"`
	Bank      *string  `json:"bank"`
	Sender    *string  `json:"sender"`
	Receiver  *string  `json:"receiver"`
	Detail    *string  `json:"detail"`
}

func (c *TransactionController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	var req UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Date != nil {
		updates["date"] = *req.Date
	}
	if req.Time != nil {
		updates["time"] = *req.Time
	}
	if req.Reference != nil {
		updates["reference"] = *req.Reference
	}
	if req.Bank != nil {
		updates["bank"] = *req.Bank
	}
	if req.Sender != nil {
		updates["sender"] = *req.Sender
	}
	if req.Receiver != nil {
		updates["receiver"] = *req.Receiver
	}
	if req.Detail != nil {
		updates["detail"] = *req.Detail
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields to update",
		})
		return
	}

	transaction, err := c.service.Update(uint(id), updates)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Transaction updated successfully",
		"transaction": transaction,
	})
}
