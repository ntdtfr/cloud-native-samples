# Cloud Native Microservices
Samples Cloud Native projects using modern technologies (Kubernetes, Docker, Java, Golang, Node.JS, ...)

## Architecture

Our System is build with the following microservices:

- User Service (Java/Spring Boot): User management and authentication
- Product Service (Golang): Product catalog and inventory
- Order Service (Node.js): Order processing and management
- Notification Service (Python): Handles emails, SMS, and other notifications
- Analytics Service (Python): Collects and processes metrics/analytics

## Project structure

```bash
cloud-native-samples/
├── .github/
│   └── workflows/           # GitHub Actions CI/CD pipelines
│       ├── user-service-ci.yml
│       ├── product-service-ci.yml
│       ├── order-service-ci.yml
│       ├── notification-service-ci.yml
│       └── analytics-service-ci.yml
├── docs
│   ├── api
│   ├── architecture
│   │   └── architecture.md
│   └── CONTRIBUTION.md
├── infrastructure
│   ├── api-gateway          # API Gateway configuration
│   │   ├── kong
│   │   ├── krakend
│   │   └── skaffold.yml
│   ├── databases            # Databases configuration
│   │   ├── duckdb
│   │   ├── mongodb
│   │   ├── postgres
│   │   ├── redis
│   │   └── skaffold.yml
│   ├── helm                 # Helm charts for deployment
│   │   └── order-service
│   ├── ingress              # Kubernetes ingress configuration
│   ├── message-brokers      # Message brokers configuration
│   │   └── rabbitmq
│   ├── monitoring           # Monitoring and logging configuration
│   ├── security             # Security configuration
│   │   ├── keycloak         
│   │   └── skaffold.yml
│   ├── skaffold.yml
│   └── terraform
├── services
│   ├── user-service         # Java/Spring Boot
│   ├── product-service      # Golang
│   ├── order-service        # Node.js
│   ├── notification-service # Python
│   └── analytics-service    # Python
├── README.md
├── docker-compose.yml
├── skaffold.yml             # Kubernetes deployment using Skaffold
└── tests
    ├── integration          # Integration tests
    ├── e2e                  # End-to-end tests with Cypress
    └── performance          # Performance testing scripts
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Skaffold v2+
- Minikube

### Clone the project

Clone the repository

```bash
git clone https://github.com/nntran/cloud-native-samples.git
cd cloud-native-samples
```

### Run with Docker Compose

* Launch all services in background
```bash
# Build and start containers
docker compose up -d

# View logs
docker compose logs -f
```

* Build and launch all services in foreground
```bash
docker compose up --build
```

* Build and launch some specific services
```bash
# Build and run only keycloak
docker compose up --build keycloak

# Build and run keycloak and order-service
docker compose up --build keycloak order-service
```

* Stop services 
```bash
docker compose down
```

### Using Skaffold to deploy to Minikube

#### 1. Start Minikube

```bash
minikube start
```

#### 2. Deploy services

* Build and deploy all services to Minikube

```bash
skaffold dev -p minikube

# or 

skaffold run -p minikube
```

* Build and deploy some specific services or groups of services to Minikube

```bash
skaffold dev -p minikube -m api-gateway,security,order-service
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
