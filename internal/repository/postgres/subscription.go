package sub_posgres_repo

import (
	"context"
	"fmt"

	"github.com/chixxx1/subscription-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepo(db *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: db,
	}
}

func (sr *SubscriptionRepo) Create(ctx context.Context, sub domain.Subscription) error {
	query := `
		INSERT INTO subscription (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := sr.db.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	if err != nil {
		return fmt.Errorf("subscriptionRepo create: %w", err)
	}
	return nil
}

func (sr *SubscriptionRepo) GetByID(ctx context.Context, id int) (*domain.Subscription, error) {
	var sub domain.Subscription
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscription
		WHERE id = $1
	`
	err := sr.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("subscriptionRepo getByID: %w", err)
	}
	return &sub, nil

}

func (sr *SubscriptionRepo) List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscription
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filter.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, filter.UserID)
		argIndex++
	}
	if filter.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, filter.ServiceName)
		argIndex++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := sr.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SubscriptionRepo.List: %w", err)
	}
	defer rows.Close()

	var subscriptions []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("subscriptionRepo list scan: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (sr *SubscriptionRepo) Update(ctx context.Context, sub domain.Subscription) error {
	query := `
		UPDATE subscription
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW()
		WHERE id = $6
	`

	_, err := sr.db.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return fmt.Errorf("SubscriptionRepo.Update: %w", err)
	}

	return nil
}

func (sr *SubscriptionRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscription WHERE id = $1`
	_, err := sr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SubscriptionRepo.Delete: %w", err)
	}

	return nil
}

func (sr *SubscriptionRepo) GetTotalCost(ctx context.Context, req domain.TotalCostRequest) (int64, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)::bigint
		FROM subscription
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if req.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, req.UserID)
		argIndex++
	}

	if req.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, req.ServiceName)
		argIndex++
	}

	query += fmt.Sprintf(" AND start_date <= $%d AND (end_date IS NULL OR end_date >= $%d)", argIndex, argIndex+1)
	args = append(args, req.PeriodEnd, req.PeriodStart)

	var total int64
	err := sr.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("SubscriptionRepo.GetTotalCost: %w", err)
	}

	return total, nil
}
