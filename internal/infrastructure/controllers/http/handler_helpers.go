package http

import (
	model "Subscription_Service/internal/domain/subscription"
	"Subscription_Service/internal/infrastructure/controllers/dto"
)

func toResponse(s model.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
