package graph

import (
	"curso-go-clean-arch/internal/container"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	container *container.Container
}

// NewResolver creates a new resolver with dependencies
func NewResolver(container *container.Container) *Resolver {
	return &Resolver{
		container: container,
	}
}
