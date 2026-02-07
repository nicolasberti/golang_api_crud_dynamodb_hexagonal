# AGENTS.md

## Project Overview

This is a Go REST API following hexagonal architecture (clean architecture) patterns with AWS DynamoDB as the database. The project implements best practices including SOLID principles, DRY, and Go-specific conventions. It's designed to be cloud-native with Terraform for Infrastructure as Code (IaC).

## Architecture

```
cmd/api/
├── main.go                    # Application entry point

internal/
├── core/
│   ├── domain/
│   │   └── product.go         # Business entities and domain logic
│   ├── ports/
│   │   ├── repositories.go    # Repository interfaces
│   │   └── services.go        # Service interfaces
│   └── services/
│       └── product_service.go # Business logic implementation
├── adapters/
│   ├── http/
│   │   └── product_handler.go # HTTP layer (primary adapter)
│   └── repository/
│       └── dynamodb.go       # DynamoDB implementation (secondary adapter)
└── platform/
    ├── config/
    │   └── config.go          # Configuration management
    └── logger/
        └── logger.go          # Logging utilities

terraform/
├── main.tf                    # AWS resources definition
├── variables.tf               # Terraform variables
└── outputs.tf                 # Output values
```

## Development Guidelines

### Code Style & Best Practices

1. **Go Conventions**
   - Use `gofmt` and `golint` for code formatting
   - Follow standard Go project layout
   - Use meaningful variable and function names
   - Export functions/types only when necessary
   - Use short variable declarations (`:=`) within functions

2. **SOLID Principles**
   - **S**: Single Responsibility - Each component has one purpose
   - **O**: Open/Closed - Extensible through interfaces, not modification
   - **L**: Liskov Substitution - Interfaces can be swapped
   - **I**: Interface Segregation - Small, focused interfaces
   - **D**: Dependency Inversion - Depend on abstractions, not concretions

3. **Clean Architecture Patterns**
   - **Domain Layer**: Pure business logic, no external dependencies
   - **Ports**: Define contracts between layers
   - **Adapters**: Implement ports for external concerns
   - **Dependency Injection**: Constructor injection, no global state

4. **Error Handling**
   - Use explicit error returns, never ignore errors
   - Wrap errors with context using `fmt.Errorf`
   - Define domain-specific errors in domain package
   - Log errors at appropriate levels (Info, Warn, Error)

5. **Testing Strategy**
   - Unit tests for business logic (services, domain)
   - Integration tests for adapters
   - Use interfaces for easy mocking
   - Table-driven tests for multiple scenarios

### AWS & Cloud-Native Guidelines

1. **DynamoDB Best Practices**
   - Use appropriate partition keys for even distribution
   - Implement retry logic with exponential backoff
   - Use batch operations when possible
   - Enable encryption at rest and point-in-time recovery

2. **Security**
   - Use IAM roles instead of access keys
   - Principle of least privilege for permissions
   - Enable VPC endpoints when applicable
   - Never commit secrets to version control

3. **Observability**
   - Structured logging with context
   - Include correlation IDs for tracing
   - Monitor key business and technical metrics
   - Set up appropriate CloudWatch alerts

## Commands

### Development Commands
```bash
# Run the application
go run cmd/api/main.go

# Build the application
go build -o bin/product-api cmd/api/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Lint code
golint ./...
```

### Terraform Commands
```bash
# Initialize Terraform
cd terraform && terraform init

# Plan infrastructure changes
terraform plan

# Apply infrastructure changes
terraform apply

# Destroy infrastructure
terraform destroy
```

### Testing & Quality Assurance
```bash
# Run all tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Check for vulnerabilities
go list -json -m all | nancy sleuth
```

## Dependencies

### Core Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/aws/aws-sdk-go-v2` - AWS SDK v2
- `github.com/google/uuid` - UUID generation
- `log/slog` - Structured logging

### Testing Dependencies
- `go.uber.org/mock` - Mocking framework

## Environment Variables

```bash
# Server Configuration
PORT=8080
LOG_LEVEL=info

# AWS Configuration
AWS_REGION=us-east-1
DYNAMODB_TABLE=products
```

## API Endpoints

```
GET    /health                 # Health check
GET    /api/v1/products        # List all products
POST   /api/v1/products        # Create new product
GET    /api/v1/products/:id    # Get product by ID
PUT    /api/v1/products/:id    # Update product
DELETE /api/v1/products/:id    # Delete product
```

## Skills Auto-Invocation

### Available Skills

#### golang-pro

This skill provides specialized Go development assistance including:
- Architecture pattern implementation
- Performance optimization
- Advanced Go idioms and patterns
- Best practices review
- Code generation and refactoring

**Usage:** Auto-invoked for complex Go tasks, architecture decisions, and advanced optimizations.

To use the golang-pro skill:
```
/skill golang-pro "optimize the DynamoDB repository for better performance"
```

### Skill Invocation Patterns

1. **Auto-invocation triggers:**
   - Complex refactoring tasks
   - Architecture pattern questions
   - Performance optimization requests
   - Advanced Go feature usage
   - Cloud-native implementation challenges

2. **Manual invocation examples:**
   ```
   /skill golang-pro "review hexagonal architecture implementation"
   /skill golang-pro "add caching layer to services"
   /skill golang-pro "implement graceful shutdown patterns"
   ```

## Monitoring & Observability

### Key Metrics to Monitor
- HTTP request latency and error rates
- DynamoDB operation latency and throttling
- Application memory and CPU usage
- Business metrics (products created/updated/deleted)

### Logging Strategy
- Use structured logging with `log/slog`
- Include correlation IDs for request tracing
- Log at appropriate levels (Error for issues, Info for important events, Debug for detailed tracing)
- Never log sensitive information (PII, credentials)

## Deployment Considerations

### Containerization
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Environment-specific Configuration
- Use environment variables for configuration
- Implement configuration validation at startup
- Support different profiles (dev, staging, prod)
- Use secrets management for sensitive data