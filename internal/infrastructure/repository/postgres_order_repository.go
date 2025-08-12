package repository

import (
	"context"
	"database/sql"
	"fmt"

	"curso-go-clean-arch/internal/domain/entity"
	"curso-go-clean-arch/internal/domain/repository"

	"github.com/google/uuid"
)

// PostgresOrderRepository implements the OrderRepository interface using PostgreSQL
type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a new instance of PostgresOrderRepository
func NewPostgresOrderRepository(db *sql.DB) repository.OrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

// Create saves a new order to the database
func (r *PostgresOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders (id, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, order.ID, order.Description, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}

	return nil
}

// List retrieves all orders from the database
func (r *PostgresOrderRepository) List(ctx context.Context) ([]*entity.Order, error) {
	query := `
		SELECT id, description, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying orders: %w", err)
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.Description, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// GetByID retrieves an order by its ID
func (r *PostgresOrderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	query := `
		SELECT id, description, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &entity.Order{}
	err = r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID, &order.Description, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	return order, nil
}

// Update updates an existing order in the database
func (r *PostgresOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	query := `
		UPDATE orders
		SET description = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, order.Description, order.UpdatedAt, order.ID)
	if err != nil {
		return fmt.Errorf("error updating order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Delete removes an order from the database
func (r *PostgresOrderRepository) Delete(ctx context.Context, id string) error {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}
