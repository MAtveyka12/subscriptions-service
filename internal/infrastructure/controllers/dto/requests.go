package dto

import (
	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required,min=2,max=100" example:"Yandex Plus"`
	Price       int        `json:"price" binding:"required,gte=0" example:"400"`
	UserID      uuid.UUID  `json:"user_id" binding:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   CustomTime `json:"start_date" binding:"required" example:"2025-07-01"`
	EndDate     CustomTime `json:"end_date,omitempty" example:"2025-12-31"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string     `json:"service_name,omitempty" binding:"omitempty,min=2,max=100"`
	Price       *int        `json:"price,omitempty" binding:"omitempty,gte=0"`
	StartDate   *CustomTime `json:"start_date,omitempty"`
	EndDate     *CustomTime `json:"end_date,omitempty"`
}
