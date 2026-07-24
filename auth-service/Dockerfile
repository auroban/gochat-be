# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Copy workspace file from root
COPY go.work go.work.sum* ./

# Copy both modules
COPY apis ./apis
COPY auth-service ./auth-service

# Download dependencies
WORKDIR /app/auth-service
RUN go mod download

# Build the application
RUN go build -o auth-service ./cmd/server

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service/auth-service .
COPY --from=builder /app/auth-service/configs/config.yaml .
COPY --from=builder /app/auth-service/migrations /app/migrations
EXPOSE 4051
CMD ["./auth-service"]
