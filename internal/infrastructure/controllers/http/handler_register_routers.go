package http

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/subscriptions", h.Create)
	r.GET("/subscriptions/:id", h.Read)
	r.PUT("/subscriptions/:id", h.Update)
	r.DELETE("/subscriptions/:id", h.Delete)
	r.GET("/subscriptions", h.List)
	r.GET("/subscriptions/cost", h.CalculateCost)
}
