package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	model "Subscription_Service/internal/domain/subscription"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *model.Subscription) error
	Read(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	Update(ctx context.Context, s *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.Subscription, error)
	FindFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) ([]model.Subscription, error)
	CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate time.Time) (int64, error)
}

type subscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (sr *subscriptionRepository) Create(ctx context.Context, s *model.Subscription) error {
	query := `
	INSERT INTO subscription (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8) 
	`
	now := time.Now().UTC()

	s.ID = uuid.New()
	s.CreatedAt = now
	s.UpdatedAt = now

	_, err := sr.db.ExecContext(ctx, query, s.ID, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create subscription")
	}

	return nil
}

func (sr *subscriptionRepository) Read(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var s model.Subscription

	err := sr.db.GetContext(ctx, &s, `SELECT * FROM subscription WHERE id=$1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription with id %s not found", id)
		}

		return nil, fmt.Errorf("failed to get subscription %s: %s", id, err.Error())
	}

	return &s, err
}

func (sr *subscriptionRepository) Update(ctx context.Context, s *model.Subscription) error {
	s.UpdatedAt = time.Now().UTC()

	result, err := sr.db.ExecContext(ctx, `UPDATE subscription SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5, updated_at=$6 WHERE id=$7`,
		s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.UpdatedAt, s.ID)

	if err != nil {
		return fmt.Errorf("failed to update subscription %s: %s", s.ID, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for subscription %s: %s", s.ID, err.Error())
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription with id %s not found", s.ID)
	}

	return nil
}

func (sr *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := sr.db.ExecContext(ctx, `DELETE FROM subscription WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription %s: %s", id, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for subscription %s: %s", id, err.Error())
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription with id %s not found", id)
	}

	return nil
}

func (sr *subscriptionRepository) List(ctx context.Context, limit, offset int) (subs []model.Subscription, err error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	subs = []model.Subscription{}
	err = sr.db.SelectContext(ctx, &subs, `SELECT * FROM subscription ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("list subscription: %s", err.Error())
	}

	return subs, err
}

func (sr *subscriptionRepository) FindFiltered(ctx context.Context, userID *uuid.UUID, serviceName *string, limit, offset int) (subs []model.Subscription, err error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	conds := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)

	if userID != nil {
		conds = append(conds, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, *userID)
	}

	if serviceName != nil && strings.TrimSpace(*serviceName) != "" {
		conds = append(conds, fmt.Sprintf("service_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+strings.TrimSpace(*serviceName)+"%")
	}

	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscription`
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	subs = []model.Subscription{}

	err = sr.db.SelectContext(ctx, &subs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find filtered subscription: %s", err.Error())
	}

	return subs, nil
}

func (sr *subscriptionRepository) CalculateCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate time.Time) (int64, error) {
	ps := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	pe := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	conds := make([]string, 0, 4)
	args := make([]interface{}, 0, 6)

	if userID != nil {
		conds = append(conds, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, *userID)
	}

	if serviceName != nil && strings.TrimSpace(*serviceName) != "" {
		conds = append(conds, fmt.Sprintf("service_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+strings.TrimSpace(*serviceName)+"%")
	}

	conds = append(conds, fmt.Sprintf("start_date <= $%d", len(args)+1))
	args = append(args, pe)
	conds = append(conds, fmt.Sprintf("(end_date IS NULL OR end_date >= $%d)", len(args)+1))
	args = append(args, ps)

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(`
	WITH filtered AS (
	SELECT price,
	GREATEST(start_date, $%d::date) AS s,
	LEAST(COALESCE(end_date, $%d::date), $%d::date) AS e
	FROM subscription
	%s
	)
	SELECT COALESCE(SUM(price * (
	(date_part('year', e)::int - date_part('year', s)::int) * 12
	+ (date_part('month', e)::int - date_part('month', s)::int) + 1
	)), 0) AS total
	FROM filtered`, len(args)+1, len(args)+2, len(args)+2, where)

	args = append(args, ps, pe)

	var total sql.NullInt64
	if err := sr.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, fmt.Errorf("calculate cost: %s", err.Error())
	}

	if !total.Valid {
		return 0, nil
	}

	return total.Int64, nil
}
