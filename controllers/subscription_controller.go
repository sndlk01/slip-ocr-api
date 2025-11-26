package controllers

import (
	"net/http"
	"ocr-api/models"
	"ocr-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	service *services.SubscriptionService
}

func NewSubscriptionController() *SubscriptionController {
	return &SubscriptionController{service: services.NewSubscriptionService()}
}

type CreateSubscriptionRequest struct {
	Name            string  `json:"name" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	Category        string  `json:"category"`
	BillingCycle    string  `json:"billing_cycle" binding:"required"`
	NextBillingDate string  `json:"next_billing_date"`
}

func (c *SubscriptionController) Create(ctx *gin.Context) {
	var req CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sub := &models.Subscription{
		Name:            req.Name,
		Amount:          req.Amount,
		Category:        req.Category,
		BillingCycle:    req.BillingCycle,
		NextBillingDate: req.NextBillingDate,
		IsActive:        true,
	}

	if err := c.service.Create(sub); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Subscription created", "subscription": sub})
}

func (c *SubscriptionController) GetAll(ctx *gin.Context) {
	subs, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"subscriptions": subs})
}

func (c *SubscriptionController) Delete(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}
