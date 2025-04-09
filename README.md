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
├── infrastructure/
│   ├── k8s/                 # Kubernetes manifests
│   │   ├── ingress/
│   │   ├── databases/
│   │   ├── services/
│   │   └── monitoring/
│   ├── helm/                # Helm charts for deployment
│   └── terraform/           # Infrastructure as Code (optional)
├── services/
│   ├── user-service/        # Java/Spring Boot
│   ├── product-service/     # Golang
│   ├── order-service/       # Node.js
│   ├── notification-service/ # Python
│   └── analytics-service/   # Python
├── api-gateway/             # API Gateway configuration
├── keycloak/                # Keycloak configuration
├── docker-compose.yml       # Local development environment
├── skaffold.yml             # Simplify Kubernetes development
└── tests/
    ├── integration/         # Integration tests
    ├── e2e/                 # End-to-end tests with Cypress
    └── performance/         # Performance testing scripts
```
