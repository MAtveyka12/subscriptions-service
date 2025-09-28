package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID          uuid.UUID `json:"id" example:"a3e7f924-7d11-4f36-91bb-8f69cb1c1a91"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserID      uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time `json:"start_date" example:"2025-07-01"`
	EndDate     time.Time `json:"end_date,omitempty" example:"2025-12-31"`
	CreatedAt   time.Time `json:"created_at" example:"2025-07-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-07-02T12:00:00Z"`
}

type CostResponse struct {
	Total int64 `json:"total" example:"1200"`
}
