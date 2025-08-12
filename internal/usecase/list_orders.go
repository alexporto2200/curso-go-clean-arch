package usecase

import (
	"context"

	"curso-go-clean-arch/internal/domain/repository"
)

// ListOrdersOutput represents the output data for listing orders
type ListOrdersOutput struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListOrdersUseCase handles the business logic for listing orders
type ListOrdersUseCase struct {
	orderRepository repository.OrderRepository
}

// NewListOrdersUseCase creates a new instance of ListOrdersUseCase
func NewListOrdersUseCase(orderRepository repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		orderRepository: orderRepository,
	}
}

// Execute performs the list orders operation
func (uc *ListOrdersUseCase) Execute(ctx context.Context) ([]*ListOrdersOutput, error) {
	// Get orders from repository
	orders, err := uc.orderRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to output format
	var output []*ListOrdersOutput
	for _, order := range orders {
		output = append(output, &ListOrdersOutput{
			ID:          order.ID.String(),
			Description: order.Description,
			CreatedAt:   order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return output, nil
}
