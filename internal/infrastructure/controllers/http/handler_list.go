package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	model "Subscription_Service/internal/domain/subscription"
	"Subscription_Service/internal/infrastructure/controllers/dto"
)

func (h *Handler) List(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		return
	}

	var userID *uuid.UUID
	if v := c.Query("user_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			userID = &id
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
	}

	var serviceName *string
	if v := c.Query("service_name"); v != "" {
		serviceName = &v
	}

	var subs []model.Subscription

	if userID != nil || serviceName != nil {
		subs, err = h.service.FindFiltered(c, userID, serviceName, limit, offset)
	} else {
		subs, err = h.service.List(c, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]dto.SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		resp = append(resp, toResponse(s))
	}

	c.JSON(http.StatusOK, resp)
}
