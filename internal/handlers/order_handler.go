package handlers

import (
	"curso-go-clean-arch/internal/container"
	"curso-go-clean-arch/internal/handlers/dto"
	"curso-go-clean-arch/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	container *container.Container
	validate  *validator.Validate
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(container *container.Container) *OrderHandler {
	return &OrderHandler{
		container: container,
		validate:  validator.New(),
	}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to use case input
	input := usecase.CreateOrderInput{
		Description: req.Description,
	}

	// Execute use case
	output, err := h.container.CreateOrderUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response
	response := &dto.OrderResponse{
		ID:          output.ID,
		Description: output.Description,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Execute use case
	output, err := h.container.ListOrdersUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, "Failed to list orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response
	var orders []*dto.OrderResponse
	for _, order := range output {
		orders = append(orders, &dto.OrderResponse{
			ID:          order.ID,
			Description: order.Description,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		})
	}

	response := &dto.ListOrdersResponse{
		Orders: orders,
		Total:  len(orders),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
