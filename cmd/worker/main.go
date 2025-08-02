package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awesomeapibrasil/gateway-worker/internal/grpc"
	"github.com/awesomeapibrasil/gateway-worker/internal/health"
	"github.com/awesomeapibrasil/gateway-worker/internal/queue"
)

const (
	defaultGRPCPort = "8080"
	defaultHTTPPort = "8081"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize services
	healthService := health.New()
	queueService := queue.New()
	grpcService := grpc.New(queueService)

	// Start gRPC server
	grpcPort := getEnv("GRPC_PORT", defaultGRPCPort)
	go func() {
		if err := startGRPCServer(ctx, grpcService, grpcPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start HTTP health server
	httpPort := getEnv("HTTP_PORT", defaultHTTPPort)
	go func() {
		if err := startHTTPServer(ctx, healthService, httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start background job processor
	go func() {
		if err := queueService.Start(ctx); err != nil {
			log.Fatalf("Failed to start queue service: %v", err)
		}
	}()

	log.Printf("Worker service started - gRPC on :%s, HTTP on :%s", grpcPort, httpPort)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker service...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	cancel() // Cancel main context
	<-shutdownCtx.Done()

	log.Println("Worker service stopped")
}

func startGRPCServer(ctx context.Context, service *grpc.Service, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	return service.Serve(ctx, lis)
}

func startHTTPServer(ctx context.Context, healthService *health.Service, port string) error {
	return healthService.Serve(ctx, ":"+port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}