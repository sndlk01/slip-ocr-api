package controllers

import (
	"net/http"
	"ocr-api/services"
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
