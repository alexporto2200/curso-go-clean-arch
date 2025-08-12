package server

import (
	"curso-go-clean-arch/internal/container"
	"curso-go-clean-arch/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// RESTServer represents the REST API server
type RESTServer struct {
	router    *mux.Router
	container *container.Container
	port      string
}

// NewRESTServer creates a new REST server
func NewRESTServer(container *container.Container) *RESTServer {
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "8081"
	}

	return &RESTServer{
		router:    mux.NewRouter(),
		container: container,
		port:      port,
	}
}

// SetupRoutes configures all the routes for the REST API
func (s *RESTServer) SetupRoutes() {
	// Create handlers
	orderHandler := handlers.NewOrderHandler(s.container)

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Orders routes
	orders := api.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("", orderHandler.ListOrders).Methods("GET")
	orders.HandleFunc("", orderHandler.CreateOrder).Methods("POST")

	// Root redirect to health
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/health", http.StatusMovedPermanently)
	})

	// Add middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
}

// healthCheck handles health check requests
func (s *RESTServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "orders-api"}`))
}

// loggingMiddleware logs all requests
func (s *RESTServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware handles CORS
func (s *RESTServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Start starts the REST server
func (s *RESTServer) Start() error {
	s.SetupRoutes()

	log.Printf("REST API server starting on port %s", s.port)
	log.Printf("Health check: http://localhost:%s/health", s.port)
	log.Printf("API base: http://localhost:%s/api/v1", s.port)

	return http.ListenAndServe(":"+s.port, s.router)
}
