package grpc

import (
	"context"
	"curso-go-clean-arch/internal/container"
	"curso-go-clean-arch/internal/usecase"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	order "curso-go-clean-arch/proto"
)

// OrderServer implements the gRPC OrderService
type OrderServer struct {
	order.UnimplementedOrderServiceServer
	container *container.Container
}

// NewOrderServer creates a new gRPC order server
func NewOrderServer(container *container.Container) *OrderServer {
	return &OrderServer{
		container: container,
	}
}

// CreateOrder implements the CreateOrder RPC method
func (s *OrderServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	// Convert to use case input
	input := usecase.CreateOrderInput{
		Description: req.Description,
	}

	// Execute use case
	output, err := s.container.CreateOrderUseCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	// Convert to protobuf response
	createdAt, _ := time.Parse(time.RFC3339, output.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, output.UpdatedAt)

	protoOrder := &order.Order{
		Id:          output.ID,
		Description: output.Description,
		CreatedAt:   timestamppb.New(createdAt),
		UpdatedAt:   timestamppb.New(updatedAt),
	}

	return &order.CreateOrderResponse{
		Order: protoOrder,
	}, nil
}

// ListOrders implements the ListOrders RPC method
func (s *OrderServer) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	// Execute use case
	output, err := s.container.ListOrdersUseCase.Execute(ctx)
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		return nil, status.Error(codes.Internal, "failed to list orders")
	}

	// Convert to protobuf response
	var protoOrders []*order.Order
	for _, orderOutput := range output {
		createdAt, _ := time.Parse(time.RFC3339, orderOutput.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, orderOutput.UpdatedAt)

		protoOrder := &order.Order{
			Id:          orderOutput.ID,
			Description: orderOutput.Description,
			CreatedAt:   timestamppb.New(createdAt),
			UpdatedAt:   timestamppb.New(updatedAt),
		}
		protoOrders = append(protoOrders, protoOrder)
	}

	return &order.ListOrdersResponse{
		Orders: protoOrders,
		Total:  int32(len(protoOrders)),
	}, nil
}

// GRPCServer represents the gRPC server
type GRPCServer struct {
	server    *grpc.Server
	container *container.Container
	port      string
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(container *container.Container) *GRPCServer {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8082"
	}

	return &GRPCServer{
		server:    grpc.NewServer(),
		container: container,
		port:      port,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	// Register the service
	orderServer := NewOrderServer(s.container)
	order.RegisterOrderServiceServer(s.server, orderServer)

	// Start listening
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	log.Printf("gRPC server starting on port %s", s.port)
	log.Printf("gRPC server: localhost:%s", s.port)

	return s.server.Serve(lis)
}

// Stop stops the gRPC server
func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
