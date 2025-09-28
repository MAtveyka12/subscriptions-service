package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	model "Subscription_Service/internal/domain/subscription"
	"Subscription_Service/internal/infrastructure/repository"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, s *model.Subscription) error
	Read(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	Update(ctx context.Context, s *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.Subscription, error)
	FindFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]model.Subscription, error)
	CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate time.Time) (int64, error)
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
	logger           *slog.Logger
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository, logger *slog.Logger) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		logger:           logger,
	}
}

func (s *subscriptionService) Create(ctx context.Context, sub *model.Subscription) error {
	s.logger.Debug("Creating subscription",
		slog.String("service_name", sub.ServiceName),
		slog.String("user_id", sub.UserID.String()),
	)

	if sub.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}

	err := s.subscriptionRepo.Create(ctx, sub)
	if err != nil {
		s.logger.Error("Failed to create subscription",
			slog.String("error", err.Error()),
			slog.String("user_id", sub.UserID.String()),
		)

		return err
	}

	s.logger.Info("Subscription created successfully",
		slog.String("subscription_id", sub.ID.String()),
	)

	return nil
}

func (s *subscriptionService) Read(ctx context.Context, id uuid.UUID) (sub *model.Subscription, err error) {
	s.logger.Debug("Fetching subscription",
		slog.String("id", id.String()),
	)

	sub, err = s.subscriptionRepo.Read(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch subscription",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) Update(ctx context.Context, sub *model.Subscription) error {
	s.logger.Debug("Updating subscription",
		slog.String("id", sub.ID.String()),
	)

	err := s.subscriptionRepo.Update(ctx, sub)
	if err != nil {
		s.logger.Error("Failed to update subscription",
			slog.String("id", sub.ID.String()),
			slog.String("error", err.Error()),
		)

		return err
	}

	s.logger.Info("Subscription updated successfully",
		slog.String("id", sub.ID.String()),
	)

	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("Deleting subscription",
		slog.String("id", id.String()),
	)

	err := s.subscriptionRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete subscription",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)

		return err
	}

	s.logger.Info("Subscription deleted successfully",
		slog.String("id", id.String()),
	)

	return nil
}

func (s *subscriptionService) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	s.logger.Debug("Listing subscriptions",
		slog.Int("limit", limit),
		slog.Int("offset", offset),
	)

	subs, err := s.subscriptionRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list subscriptions",
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	s.logger.Info("Subscription listed successfully")

	return subs, nil
}

func (s *subscriptionService) FindFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]model.Subscription, error) {
	s.logger.Debug("Filtering subscriptions",
		slog.String("user_id", safeUUID(userID)),
		slog.String("service_name", safeStr(serviceName)),
	)

	subs, err := s.subscriptionRepo.FindFiltered(ctx, userID, serviceName, limit, offset)
	if err != nil {
		s.logger.Error("Failed to filter subscriptions",
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	s.logger.Info("Subscription filtered successfully",
		slog.String("user_id", safeUUID(userID)),
		slog.String("service_name", safeStr(serviceName)),
	)

	return subs, nil
}

func (s *subscriptionService) CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate time.Time) (int64, error) {
	s.logger.Debug("Calculating subscription cost",
		slog.String("user_id", safeUUID(userID)),
		slog.String("service_name", safeStr(serviceName)),
		slog.Time("start_date", startDate),
		slog.Time("end_date", endDate),
	)

	if endDate.Before(startDate) {
		return 0, fmt.Errorf("endDate cannot be before startDate")
	}

	total, err := s.subscriptionRepo.CalculateCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to calculate cost",
			slog.String("error", err.Error()),
		)

		return 0, err
	}

	s.logger.Info("Subscription calculated successfully",
		slog.String("user_id", safeUUID(userID)),
		slog.String("service_name", safeStr(serviceName)),
	)

	return total, nil
}

func safeUUID(u *uuid.UUID) string {
	if u == nil {
		return ""
	}

	return u.String()
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
