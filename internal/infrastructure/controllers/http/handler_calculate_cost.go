package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"Subscription_Service/internal/infrastructure/controllers/dto"
)

func (h *Handler) CalculateCost(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date required"})
		return
	}

	ps, err := time.Parse("2006-01", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected YYYY-MM"})
		return
	}

	pe, err := time.Parse("2006-01", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected YYYY-MM"})
		return
	}

	var userID *uuid.UUID
	if v := c.Query("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			userID = &id
		}
	}

	var serviceName *string
	if v := c.Query("service_name"); v != "" {
		serviceName = &v
	}

	total, err := h.service.CalculateCost(c, userID, serviceName, ps, pe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CostResponse{Total: total})
}
