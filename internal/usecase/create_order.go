package usecase

import (
	"context"

	"curso-go-clean-arch/internal/domain/entity"
	"curso-go-clean-arch/internal/domain/repository"
)

// CreateOrderInput represents the input data for creating an order
type CreateOrderInput struct {
	Description string `json:"description" validate:"required"`
}

// CreateOrderOutput represents the output data for creating an order
type CreateOrderOutput struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateOrderUseCase handles the business logic for creating orders
type CreateOrderUseCase struct {
	orderRepository repository.OrderRepository
}

// NewCreateOrderUseCase creates a new instance of CreateOrderUseCase
func NewCreateOrderUseCase(orderRepository repository.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepository: orderRepository,
	}
}

// Execute performs the create order operation
func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
	// Create new order entity
	order := entity.NewOrder(input.Description)

	// Save to repository
	if err := uc.orderRepository.Create(ctx, order); err != nil {
		return nil, err
	}

	// Return output
	return &CreateOrderOutput{
		ID:          order.ID.String(),
		Description: order.Description,
		CreatedAt:   order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
