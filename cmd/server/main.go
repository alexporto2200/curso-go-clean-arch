package main

import (
	"context"
	"curso-go-clean-arch/graph"
	"curso-go-clean-arch/internal/container"
	"curso-go-clean-arch/internal/grpc"
	"curso-go-clean-arch/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	defaultGraphQLPort = "8080"
	defaultRESTPort    = "8081"
)

func main() {
	// Initialize container with dependencies
	container, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Create GraphQL server
	graphQLServer := createGraphQLServer(container)

	// Create REST server
	restServer := server.NewRESTServer(container)

	// Create gRPC server
	grpcServer := grpc.NewGRPCServer(container)

	// Start servers in goroutines
	go func() {
		port := os.Getenv("GRAPHQL_PORT")
		if port == "" {
			port = defaultGraphQLPort
		}

		log.Printf("GraphQL server starting on port %s", port)
		log.Printf("GraphQL playground: http://localhost:%s/", port)

		if err := http.ListenAndServe(":"+port, graphQLServer); err != nil {
			log.Fatalf("GraphQL server failed: %v", err)
		}
	}()

	go func() {
		if err := restServer.Start(); err != nil {
			log.Fatalf("REST server failed: %v", err)
		}
	}()

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Here you could add graceful shutdown logic for both servers
	log.Println("Servers stopped")
}

func createGraphQLServer(container *container.Container) http.Handler {
	// Create resolver with dependencies
	resolver := graph.NewResolver(container)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Create mux for GraphQL
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	return mux
}
