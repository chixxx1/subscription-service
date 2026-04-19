package sub_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/chixxx1/subscription-service/internal/domain"
	"go.uber.org/zap"
)

var (
	ErrInvalidPrice         = errors.New("price must be positive")
	ErrInvalidDates         = errors.New("end date cannot be before start date")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type SubscriptionService struct {
	repo   domain.SubscriptionRepository
	logger *zap.Logger
}

func NewSubscriptionService(repo domain.SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, sub domain.Subscription) error {
	s.logger.Info("creating subscription", zap.String("service_name", sub.ServiceName), zap.String("user_id", sub.UserID))

	if sub.Price <= 0 {
		s.logger.Warn("invalid price", zap.Int("price", sub.Price))
		return ErrInvalidPrice
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		s.logger.Warn("invalid dates", zap.Time("start", sub.StartDate), zap.Time("end", *sub.EndDate))
		return ErrInvalidDates
	}

	err := s.repo.Create(ctx, sub)
	if err != nil {
		s.logger.Error("failed to create subscription in db", zap.Error(err))
		return fmt.Errorf("service create: %w", err)
	}

	s.logger.Info("subscription created successfully", zap.Int("id", sub.ID))
	return nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int) (*domain.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service get by id: %w", err)
	}
	return sub, nil
}

func (s *SubscriptionService) List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50 
	}
	
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	
	if filter.Limit > 1000 {
		filter.Limit = 1000
	}

	subs, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service list: %w", err)
	}
	return subs, nil
}

func (s *SubscriptionService) Update(ctx context.Context, sub domain.Subscription) error {
	s.logger.Info("updating subscription", zap.Int("id", sub.ID))

	if sub.Price <= 0 {
		return ErrInvalidPrice
	}
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return ErrInvalidDates
	}

	err := s.repo.Update(ctx, sub)
	if err != nil {
		return fmt.Errorf("service update: %w", err)
	}
	return nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int) error {
	s.logger.Info("deleting subscription", zap.Int("id", id))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service delete: %w", err)
	}
	return nil
}

func (s *SubscriptionService) GetTotalCost(ctx context.Context, req domain.TotalCostRequest) (int64, error) {
	s.logger.Info("calculating total cost",
		zap.String("user_id", req.UserID),
		zap.String("service", req.ServiceName),
		zap.Time("period_start", req.PeriodStart),
		zap.Time("period_end", req.PeriodEnd),
	)

	total, err := s.repo.GetTotalCost(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("service get total cost: %w", err)
	}

	s.logger.Info("total cost calculated", zap.Int64("total", total))
	return total, nil
}
