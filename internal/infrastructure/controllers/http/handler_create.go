package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	model "Subscription_Service/internal/domain/subscription"
	"Subscription_Service/internal/infrastructure/controllers/dto"
)

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub := model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate.Time,
	}

	if !req.EndDate.Time.IsZero() {
		sub.EndDate = req.EndDate.Time
	}

	if err := h.service.Create(c, &sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toResponse(sub))
}
