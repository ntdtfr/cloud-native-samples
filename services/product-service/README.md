# Product Service

## Overview


## Architecture


## Project structure

```
product-service/
├── api/                    # HTTP API handlers and middleware
│   ├── handlers/
│   ├── middleware/
│   ├── router.go           # Gin router setup
│   └── swagger/            # OpenAPI/Swagger documentation
├── cmd/                    # Command line entry points
│   └── server/
│       └── main.go
├── config/                 # Configuration management
│   ├── config.go
│   └── config.yaml
├── deployments/            # Deployment configurations
│   └── kubernetes/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── ingress.yaml
├── docs/                   # Documentation
├── internal/               # Private application code
│   ├── domain/             # Domain models
│   ├── repository/         # Data storage interfaces
│   └── service/            # Business logic
├── pkg/                    # Public libraries
│   ├── cache/              # Redis client
│   ├── database/           # MongoDB client
│   ├── logger/             # Logging utility
│   ├── messaging/          # RabbitMQ client
│   └── validator/          # Validation utilities
├── scripts/                # Utility scripts
├── test/                   # Test files
├── .dockerignore
├── .gitignore
├── Dockerfile              # Deployment configurations
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile                # Build automation
└── README.md
```
