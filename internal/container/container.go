package container

import (
	"database/sql"

	"curso-go-clean-arch/internal/database"
	"curso-go-clean-arch/internal/domain/repository"
	postgres "curso-go-clean-arch/internal/infrastructure/repository"
	"curso-go-clean-arch/internal/usecase"
)

// Container holds all dependencies
type Container struct {
	DB                 *sql.DB
	OrderRepository    repository.OrderRepository
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

// NewContainer creates and configures all dependencies
func NewContainer() (*Container, error) {
	// Database connection
	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		return nil, err
	}

	// Repository
	orderRepository := postgres.NewPostgresOrderRepository(db)

	// Use cases
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	return &Container{
		DB:                 db,
		OrderRepository:    orderRepository,
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}, nil
}

// Close closes all resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
