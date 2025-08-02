package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Service provides health checking functionality
type Service struct {
	startTime time.Time
	server    *http.Server
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string    `json:"status"`
	Uptime  string    `json:"uptime"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

// New creates a new health service
func New() *Service {
	return &Service{
		startTime: time.Now(),
	}
}

// Serve starts the HTTP health server
func (s *Service) Serve(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ready", s.readinessHandler)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error but don't exit
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

func (s *Service) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Uptime:  time.Since(s.startTime).String(),
		Version: "1.0.0", // TODO: Get from build
		Time:    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Service) readinessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if all required services are ready
	// For now, just return ready
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}