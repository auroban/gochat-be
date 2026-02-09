# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Copy workspace file from root
COPY go.work go.work.sum* ./

# Copy both modules
COPY apis ./apis
COPY user-service ./user-service

# Download dependencies
WORKDIR /app/user-service
RUN go mod download

# Build the application
RUN go build -o user-service ./cmd/server

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service/user-service .
COPY --from=builder /app/user-service/configs/config.yaml .
COPY --from=builder /app/user-service/migrations /app/migrations
EXPOSE 4050
CMD ["./user-service"]