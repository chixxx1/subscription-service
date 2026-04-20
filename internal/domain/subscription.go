package domain

import (
	"errors"
	"time"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string //UUID
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubscriptionFilter struct {
	UserID      string
	ServiceName string
	Limit       int
	Offset      int
}

type TotalCostRequest struct {
	UserID      string
	ServiceName string
	PeriodStart time.Time
	PeriodEnd   time.Time
}
