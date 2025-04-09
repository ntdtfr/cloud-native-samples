# Cloud Native Microservices
Samples Cloud Native projects using modern technologies (Kubernetes, Docker, Java, Golang, Node.JS, ...)


## Project structure

``̀ 
cloud-native-samples/
├── .github/
│   └── workflows/           # GitHub Actions CI/CD pipelines
│       ├── user-service.yml
│       ├── product-service.yml
│       ├── order-service.yml
│       ├── notification-service.yml
│       └── analytics-service.yml
├── infrastructure/
│   ├── docker-compose.yml   # Local development environment
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
├── api-gateway/             # Kong API Gateway configuration
├── keycloak/                # Keycloak configuration
└── tests/
    ├── integration/         # Integration tests
    ├── e2e/                 # End-to-end tests with Cypress
    └── performance/         # Performance testing scripts
```
