package http

import (
	"Subscription_Service/internal/application/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(serv service.Service) *Handler {
	return &Handler{
		service: serv,
	}
}
