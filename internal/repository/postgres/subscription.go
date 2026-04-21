package sub_posgres_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/chixxx1/subscription-service/internal/domain"
	"github.com/jackc/pgx/v5"
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

func (sr *SubscriptionRepo) Create(ctx context.Context, sub domain.Subscription) (int, error) {
	query := `
		INSERT INTO subscription (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var newID int
	err := sr.db.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("subscriptionRepo create: %w", err)
	}

	return newID, nil
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
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

	tag, err := sr.db.Exec(ctx, query,
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
	if tag.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

func (sr *SubscriptionRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscription WHERE id = $1`
	tag, err := sr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("SubscriptionRepo.Delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

func (sr *SubscriptionRepo) GetTotalCost(ctx context.Context, req domain.TotalCostRequest) (int64, error) {
	userID := req.UserID
	if userID == "" {
		userID = "" 
	}

	query := `
		SELECT COALESCE(SUM(price * months_cnt), 0)::bigint
		FROM (
			SELECT 
				price,
				EXTRACT(YEAR FROM age_months) * 12 + EXTRACT(MONTH FROM age_months) + 1 as months_cnt
			FROM (
				SELECT 
					price,
					AGE(
						DATE_TRUNC('month', LEAST(COALESCE(end_date, 'infinity'::date), $2::date)),
						DATE_TRUNC('month', GREATEST(start_date, $1::date))
					) as age_months
				FROM subscription
				WHERE
					($3 = '' OR user_id = $3::uuid) AND
					($4 = '' OR service_name = $4) AND
					start_date <= $2::date AND 
					(COALESCE(end_date, 'infinity'::date) >= $1::date)
			) AS inner_sub
		) AS final_sub
		WHERE months_cnt > 0
	`

	var total int64
	err := sr.db.QueryRow(ctx, query, req.PeriodStart, req.PeriodEnd, req.UserID, req.ServiceName).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("SubscriptionRepo.GetTotalCost: %w", err)
	}

	return total, nil
}
