package dto

import (
	"curso-go-clean-arch/internal/domain/entity"
	"time"
)

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	Description string `json:"description" validate:"required"`
}

// OrderResponse represents the response body for order operations
type OrderResponse struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListOrdersResponse represents the response body for listing orders
type ListOrdersResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Total  int              `json:"total"`
}

// ToEntity converts CreateOrderRequest to domain entity
func (r *CreateOrderRequest) ToEntity() *entity.Order {
	return entity.NewOrder(r.Description)
}

// FromEntity converts domain entity to OrderResponse
func FromEntity(order *entity.Order) *OrderResponse {
	return &OrderResponse{
		ID:          order.ID.String(),
		Description: order.Description,
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
	}
}

// FromEntities converts slice of domain entities to ListOrdersResponse
func FromEntities(orders []*entity.Order) *ListOrdersResponse {
	var responses []*OrderResponse
	for _, order := range orders {
		responses = append(responses, FromEntity(order))
	}

	return &ListOrdersResponse{
		Orders: responses,
		Total:  len(responses),
	}
}
