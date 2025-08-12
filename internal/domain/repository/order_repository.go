package repository

import (
	"context"

	"curso-go-clean-arch/internal/domain/entity"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	List(ctx context.Context) ([]*entity.Order, error)
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id string) error
}
