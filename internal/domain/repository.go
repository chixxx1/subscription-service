package domain

import "context"

type SubscriptionRepository interface {
	Create(ctx context.Context, sub Subscription) (int, error)
	GetByID(ctx context.Context, id int) (*Subscription, error)
	List(ctx context.Context, filter SubscriptionFilter) ([]Subscription, error)
	Update(ctx context.Context, sub Subscription) error
	Delete(ctx context.Context, id int) error
	GetTotalCost(ctx context.Context, req TotalCostRequest) (int64, error)
}
