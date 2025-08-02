# Gateway Worker Service

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-1.x-green.svg)](https://grpc.io/)

## Overview

The Gateway Worker Service is a dedicated microservice designed to handle asynchronous tasks and system management for the [Awesome API Brasil Gateway](https://github.com/awesomeapibrasil/gateway). This service is part of a two-service architecture that separates real-time proxy operations from background processing tasks.

## Architecture Reference

⚠️ **IMPORTANT**: All technical specifications, implementation details, and architectural decisions for this Worker service are documented in the [`WORKER-PURPOSE.md`](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) file in the main Gateway repository.

This document contains:
- Complete architecture overview and service separation rationale
- Detailed communication protocols (gRPC/gRPCS) between Gateway and Worker
- Certificate management flows, including temporary certificate handling
- Module specifications for configuration management, log processing, analytics, and integrations
- Migration plan with phases and success criteria
- Security considerations and network communication patterns

## Service Responsibilities

Based on the [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) specification, this Worker service handles:

### 🔐 Certificate Management
- Automatic certificate renewal via ACME protocol (Let's Encrypt, etc.)
- Certificate validation and health monitoring
- Certificate distribution to Gateway instances via gRPC
- Temporary certificate generation and deployment for failover scenarios
- Certificate backup and recovery operations

### ⚙️ Configuration Management
- WAF rule updates and compilation
- Routing configuration management
- Backend health check configuration
- Security policy updates and validation
- Configuration version control and rollback capabilities

### 📊 Log Processing & Analytics
- Access log aggregation and processing
- Security event correlation and analysis
- Audit trail generation
- Performance metrics aggregation
- Traffic pattern analysis and reporting

### 🗄️ Database Operations
- Database schema migrations
- Data cleanup and archival processes
- Backup operations and scheduling
- Performance optimization tasks

### 🔗 Integration Tasks
- External API integrations
- Third-party security feed processing
- Notification and alerting systems
- Report generation and distribution

## Technology Stack

- **Language**: Go 1.24+
- **Communication**: gRPC/gRPCS for Gateway-Worker communication
- **Queue System**: Built-in job queue with configurable workers
- **Configuration**: Environment variables and configuration files
- **Logging**: Structured logging with context
- **Health Checks**: HTTP endpoints for health and readiness probes

## Project Structure

```
.
├── cmd/worker/              # Main application entry point
├── internal/                # Private application packages
│   ├── certificate/         # Certificate management module
│   ├── config/             # Configuration management module
│   ├── log/                # Log processing module
│   ├── analytics/          # Analytics and monitoring module
│   ├── database/           # Database operations module
│   ├── integration/        # External integrations module
│   ├── grpc/               # gRPC server implementation
│   ├── health/             # Health check service
│   └── queue/              # Job queue system
├── pkg/                    # Public packages
│   ├── proto/              # Generated protobuf code
│   └── types/              # Shared type definitions
├── api/                    # API definitions (protobuf)
├── scripts/                # Build and deployment scripts
└── docs/                   # Additional documentation
```

## Quick Start

### Prerequisites

- Go 1.24 or later
- Protocol Buffers compiler (protoc)
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/awesomeapibrasil/gateway-worker.git
cd gateway-worker
```

2. Install dependencies:
```bash
go mod tidy
```

3. Generate protobuf code:
```bash
# TODO: Add protobuf generation commands
# make proto
```

4. Build the application:
```bash
go build -o bin/worker cmd/worker/main.go
```

5. Run the service:
```bash
./bin/worker
```

### Configuration

The service can be configured using environment variables:

- `GRPC_PORT`: gRPC server port (default: 8080)
- `HTTP_PORT`: HTTP health server port (default: 8081)
- `LOG_LEVEL`: Logging level (default: info)
- `QUEUE_WORKERS`: Number of job queue workers (default: 5)

### Docker

```bash
# TODO: Add Dockerfile and docker-compose.yml
docker build -t gateway-worker .
docker run -p 8080:8080 -p 8081:8081 gateway-worker
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Generation

```bash
# Generate protobuf code
make proto

# Generate mocks
make mocks
```

### Development with Gateway

This Worker service is designed to work with the main Gateway service. For complete development setup:

1. Set up the [Gateway service](https://github.com/awesomeapibrasil/gateway)
2. Configure gRPC communication between services
3. Refer to [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) for integration details

## Migration Status

This project implements the Worker service as described in the migration plan from [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md).

### Phase 1: Foundation ✅
- [x] Worker service bootstrap
- [x] Basic service framework (HTTP API, health checks)
- [x] Basic job queue system
- [x] gRPC communication infrastructure setup

### Phase 2: Certificate Management 🚧
- [ ] ACME client implementation
- [ ] Certificate storage and distribution system
- [ ] Temporary certificate generation
- [ ] Gateway integration for certificate updates

### Phase 3: Configuration Management 📋
- [ ] WAF rule management
- [ ] Configuration validation and testing
- [ ] Hot configuration updates
- [ ] Administrative interface

### Phase 4: Log Processing 📋
- [ ] Log aggregation pipeline
- [ ] Analytics and reporting engine
- [ ] Real-time analytics
- [ ] Log archival and cleanup

### Phase 5: Production Deployment 📋
- [ ] Production rollout
- [ ] Performance optimization
- [ ] Documentation and training
- [ ] Monitoring and alerting

## API Documentation

The Worker service exposes gRPC APIs for communication with the Gateway service. The API is defined in `api/gateway_worker.proto` and includes:

- Certificate management operations
- Configuration updates
- Health and status monitoring
- Job management

For detailed API documentation, see the protobuf definitions and refer to [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md).

## Health Checks

The service provides HTTP endpoints for health monitoring:

- `GET /health`: Service health status
- `GET /ready`: Service readiness status

## Contributing

1. Read the [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) specification
2. Fork the repository
3. Create a feature branch
4. Make your changes
5. Add tests
6. Submit a pull request

All contributions must align with the specifications in [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Gateway](https://github.com/awesomeapibrasil/gateway) - Main API Gateway service
- [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) - Complete architectural specification

---

> **Note**: This Worker service is designed according to the specifications in [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md). All technical discussions, implementations, and future revisions should reference this document as the authoritative source.