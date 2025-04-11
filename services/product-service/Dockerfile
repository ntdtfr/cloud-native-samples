# deployments/docker/Dockerfile

ARG GO_VERSION=1.23.3

## Stage: Build
FROM golang:${GO_VERSION}-alpine AS build

# Set workdir
WORKDIR /opt/build

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o product-service ./cmd/server/main.go

## 
## Stage: Create development image
FROM alpine:3.18 AS development

RUN apk --no-cache add ca-certificates && \
    addgroup -S ntdt && adduser -S ntran -G ntdt

# Set workdir
WORKDIR /opt/app

# Copy binary and config files from build stage to final image
COPY --chown=ntran:ntdt --from=build /opt/build/product-service .
COPY --chown=ntran:ntdt --from=build /opt/build/config/config.yaml ./config/

# Run as non-root
USER admin

EXPOSE 8080

CMD ["./product-service"]

## 
## Stage: Create production image
FROM alpine:3.18 AS production

RUN apk --no-cache add ca-certificates && \
    addgroup -S ntdt && adduser -S ntran -G ntdt

# Set workdir
WORKDIR /opt/app

# Copy binary and config files from build stage to final image
COPY --chown=ntran:ntdt --from=build /opt/build/product-service .
COPY --chown=ntran:ntdt --from=build /opt/build/config/config.yaml ./config/

# Run as non-root
USER admin

EXPOSE 8080

CMD ["./product-service"]
