package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/awesomeapibrasil/gateway-worker/internal/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Service provides gRPC functionality for Gateway-Worker communication
type Service struct {
	server      *grpc.Server
	queueService *queue.Service
}

// New creates a new gRPC service
func New(queueService *queue.Service) *Service {
	// TODO: Configure TLS credentials for secure communication
	// For now, we'll use insecure connection for development
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(10 * time.Second),
	}

	// In production, add TLS credentials:
	// creds := credentials.NewTLS(&tls.Config{...})
	// opts = append(opts, grpc.Creds(creds))

	server := grpc.NewServer(opts...)

	service := &Service{
		server:      server,
		queueService: queueService,
	}

	// TODO: Register the GatewayWorkerService
	// pb.RegisterGatewayWorkerServiceServer(server, service)

	return service
}

// Serve starts the gRPC server
func (s *Service) Serve(ctx context.Context, listener net.Listener) error {
	// Start server in a goroutine
	go func() {
		if err := s.server.Serve(listener); err != nil {
			// Log error but don't exit
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	s.server.GracefulStop()
	return nil
}

// TODO: Implement gRPC service methods when protobuf is generated
// This will include methods like:
// - UpdateCertificate
// - GetCertificateStatus  
// - DeployTemporaryCertificate
// - UpdateConfiguration
// - GetConfiguration
// - UpdateWAFRules
// - HealthCheck
// - GetWorkerStatus
// - SubmitJob
// - GetJobStatus

// getTLSCredentials configures TLS for secure gRPC communication
func (s *Service) getTLSCredentials() (credentials.TransportCredentials, error) {
	// TODO: Implement proper TLS configuration based on WORKER-PURPOSE.md
	// This should include:
	// - mTLS authentication for service-to-service communication
	// - Certificate validation
	// - Proper certificate authority setup
	
	config := &tls.Config{
		// Configure TLS settings
		MinVersion: tls.VersionTLS12,
	}
	
	return credentials.NewTLS(config), nil
}