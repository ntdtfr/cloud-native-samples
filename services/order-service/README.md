# Order Service Microservice

This is a microservice for managing orders in an e-commerce platform. The service is built with Node.js, Express, and MongoDB, following modern microservice architecture patterns.

## Features

- Create, read, update, and cancel orders
- Authentication and authorization with JWT
- Input validation
- API documentation with OpenAPI/Swagger
- Containerization with Docker
- Kubernetes deployment configuration
- Helm chart for easy deployment
- Comprehensive test suite
- Structured logging
- Error handling
- Rate limiting and security measures

## Tech Stack

- **Backend**: Node.js, Express.js
- **Database**: MongoDB
- **API Documentation**: OpenAPI 3.0, Swagger UI
- **Authentication**: JWT
- **Validation**: Joi
- **Testing**: Jest, Supertest
- **Logging**: Winston
- **Containerization**: Docker
- **Orchestration**: Kubernetes, Helm
- **Development Tools**: ESLint, Prettier, Nodemon

## Project Structure

The project follows a clean, modular architecture:

```
order-service/
├── src/               # Application source code
│   ├── api/           # API-related code
│   ├── config/        # Configuration files
│   ├── models/        # Database models
│   ├── services/      # Business logic
│   ├── utils/         # Utilities
│   ├── app.js         # Express application setup
│   └── server.js      # Server entry point
├── tests/             # Test files
├── deployment/        # Deployment configurations
├── docs/              # Documentation
└── ...                # Config files, etc.
```

## Getting Started

### Prerequisites

- Node.js 14+
- MongoDB
- Docker (optional)
- Kubernetes/Helm (optional for deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/order-service.git
   cd order-service
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Start the development server:
   ```bash
   npm run dev
   ```

### Running with Docker

```bash
# Build and run with Docker
npm run docker:build
npm run docker:run

# Or use Docker Compose
npm run docker:compose
```

### Running Tests

```bash
# Run all tests
npm test

# Run tests with coverage
npm run test:coverage
```

## API Documentation

The API is documented using OpenAPI/Swagger. After starting the server, access the documentation at:

```
http://localhost:8080/docs
```

## Deployment

### Docker Deployment

```bash
docker build -t order-service .
docker run -p 8080:3000 --env-file .env order-service
```

### Kubernetes Deployment

```bash
kubectl apply -f deployment/kubernetes/
```

### Helm Deployment

```bash
helm install order-service deployment/helm/order-service
```

## Architecture Overview

The service follows a layered architecture:

1. **API Layer**: HTTP request handling, validation, and response formatting
2. **Service Layer**: Business logic and operations
3. **Data Layer**: MongoDB models and database interactions

## Security Measures

- JWT authentication
- Helmet for HTTP headers security
- CORS configuration
- Rate limiting
- Input validation
- Error handling without exposing internals

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
