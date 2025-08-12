package entity

import (
	"time"

	"github.com/google/uuid"
)

// Order represents the order entity in the domain
type Order struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewOrder creates a new order with the given description
func NewOrder(description string) *Order {
	now := time.Now()
	return &Order{
		ID:          uuid.New(),
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateDescription updates the order description and sets the updated_at timestamp
func (o *Order) UpdateDescription(description string) {
	o.Description = description
	o.UpdatedAt = time.Now()
}
