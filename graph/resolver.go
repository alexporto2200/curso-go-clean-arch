package graph

import (
	"sync"

	"curso-go-clean-arch/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	mu     sync.Mutex
	orders []*model.Order
}
