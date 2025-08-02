# Gateway Worker Service - Getting Started Guide

This guide helps you get the Gateway Worker Service up and running based on the [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) specifications.

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)
- Protocol Buffers compiler (protoc) - for gRPC development

### 1. Clone and Setup

```bash
git clone https://github.com/awesomeapibrasil/gateway-worker.git
cd gateway-worker
```

### 2. Start with Docker Compose

The easiest way to run the complete stack:

```bash
# Start all services (Worker, PostgreSQL, Redis, Prometheus, Grafana)
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f worker
```

This will start:
- **Worker Service**: `:8080` (gRPC), `:8081` (HTTP health)
- **PostgreSQL**: Database with initialized schema
- **Redis**: Cache and session storage
- **Prometheus**: Metrics collection (`:9091`)
- **Grafana**: Monitoring dashboard (`:3000`, admin/admin)
- **Jaeger**: Distributed tracing (`:16686`)

### 3. Verify Health

```bash
# Check worker health
curl http://localhost:8081/health

# Check readiness
curl http://localhost:8081/ready

# View metrics (when implemented)
curl http://localhost:9090/metrics
```

### 4. Local Development

For Go development without Docker:

```bash
# Install dependencies
go mod download

# Build the application
make build

# Run locally
make run

# Or with live reload (requires 'air')
make dev
```

## Architecture Overview

The Worker Service implements the architecture described in [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md):

```
┌─────────────────┐    gRPC/gRPCS    ┌─────────────────┐
│  Gateway Service│◄─────────────────►│  Worker Service │
│                 │                  │                 │
│ - Proxy         │                  │ - Certificates  │
│ - Load Balancer │                  │ - Configuration │
│ - WAF           │                  │ - Log Processing│
│ - Rate Limiting │                  │ - Analytics     │
│ - Auth          │                  │ - Database Ops  │
└─────────────────┘                  │ - Integrations  │
                                     └─────────────────┘
```

## Core Modules

### 🔐 Certificate Management
- **Location**: `internal/certificate/`
- **Purpose**: ACME certificate renewal, temporary certificates, distribution
- **gRPC APIs**: `UpdateCertificate`, `GetCertificateStatus`, `DeployTemporaryCertificate`

### ⚙️ Configuration Management
- **Location**: `internal/config/`
- **Purpose**: WAF rules, routing config, security policies
- **gRPC APIs**: `UpdateConfiguration`, `UpdateWAFRules`

### 📊 Log Processing & Analytics
- **Location**: `internal/log/`, `internal/analytics/`
- **Purpose**: Log aggregation, traffic analysis, threat detection
- **Features**: Real-time processing, pattern detection, reporting

### 🗄️ Database Operations
- **Location**: `internal/database/`
- **Purpose**: Migrations, cleanup, archival, performance optimization
- **Database**: PostgreSQL with organized schemas

### 🔗 Integration Tasks
- **Location**: `internal/integration/`
- **Purpose**: External APIs, notifications, security feeds, reports
- **Features**: Webhook support, scheduled reports, threat intelligence

## Development Workflow

### 1. Making Changes

```bash
# Create feature branch
git checkout -b feature/certificate-acme

# Make changes
# ... edit files ...

# Test changes
make test
make lint

# Build and test
make build
make run
```

### 2. Working with gRPC

```bash
# Generate protobuf code after API changes
make proto

# This generates:
# - pkg/proto/*.pb.go
# - pkg/proto/*_grpc.pb.go
```

### 3. Database Changes

```bash
# Database runs automatically with docker-compose
# Schema: scripts/init-db.sql
# Access: postgresql://worker:worker_pass@localhost:5432/worker_db

# Connect to database
docker-compose exec postgres psql -U worker -d worker_db
```

### 4. Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Lint code
make lint

# Format code
make fmt
```

## Configuration

The service supports multiple configuration methods:

### Environment Variables
```bash
export GRPC_PORT=8080
export HTTP_PORT=8081
export LOG_LEVEL=info
export QUEUE_WORKERS=5
export DATABASE_URL="postgres://..."
```

### Configuration File
```bash
# config/worker.yaml
server:
  grpc_port: 8080
  http_port: 8081

queue:
  workers: 5
  buffer_size: 1000
```

## Monitoring

### Health Checks
- **Health**: `GET /health` - Service health status
- **Ready**: `GET /ready` - Service readiness
- **Metrics**: `GET /metrics` - Prometheus metrics (when implemented)

### Logging
```bash
# View structured logs
docker-compose logs -f worker

# Log levels: debug, info, warn, error
```

### Monitoring Stack
- **Prometheus**: Metrics collection and alerting
- **Grafana**: Dashboards and visualization
- **Jaeger**: Distributed tracing (optional)

## Integration with Gateway

The Worker Service is designed to work with the [Gateway Service](https://github.com/awesomeapibrasil/gateway):

### Communication Protocol
- **Transport**: gRPC with TLS (gRPCS)
- **Authentication**: Mutual TLS (mTLS)
- **Discovery**: Static configuration (expandable to Consul/Kubernetes)

### Key Integration Points
1. **Certificate Updates**: Worker → Gateway certificate deployment
2. **Configuration Changes**: Worker → Gateway config updates
3. **Health Monitoring**: Bidirectional health checks
4. **Log Streaming**: Gateway → Worker log forwarding

## Migration from Monolithic Gateway

Following the [WORKER-PURPOSE.md migration plan](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md):

### Current Status: Phase 1 ✅
- [x] Worker service foundation
- [x] Basic communication infrastructure
- [x] Database schema and job queue

### Next Steps: Phase 2 (Certificate Management)
- [ ] Implement ACME client
- [ ] Build certificate distribution
- [ ] Gateway integration

See [MIGRATION-CHECKLIST.md](docs/MIGRATION-CHECKLIST.md) for detailed progress.

## Troubleshooting

### Common Issues

**Service won't start:**
```bash
# Check logs
docker-compose logs worker

# Verify dependencies
docker-compose ps
```

**Database connection issues:**
```bash
# Verify database is running
docker-compose exec postgres pg_isready -U worker

# Check connection string
docker-compose exec worker env | grep DATABASE
```

**gRPC communication issues:**
```bash
# Test gRPC health (requires grpcurl)
grpcurl -plaintext localhost:8080 grpc.health.v1.Health/Check
```

### Logs Analysis
```bash
# Follow all logs
docker-compose logs -f

# Worker only
docker-compose logs -f worker

# Database only
docker-compose logs -f postgres
```

## Resources

- **Architecture Reference**: [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md)
- **Migration Plan**: [MIGRATION-CHECKLIST.md](docs/MIGRATION-CHECKLIST.md)
- **API Documentation**: [gateway_worker.proto](api/gateway_worker.proto)
- **Main Gateway**: [awesomeapibrasil/gateway](https://github.com/awesomeapibrasil/gateway)

## Support

For questions and support:
1. Review [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md)
2. Check existing issues
3. Create new issue with detailed description
4. Reference architectural specification in discussions

---

> **Important**: All technical decisions and implementations must align with [WORKER-PURPOSE.md](https://github.com/awesomeapibrasil/gateway/blob/main/WORKER-PURPOSE.md) specifications.