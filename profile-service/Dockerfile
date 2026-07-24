# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Copy workspace file from root
COPY go.work go.work.sum* ./

# Copy both modules
COPY apis ./apis
COPY profile-service ./profile-service

# Download dependencies
WORKDIR /app/profile-service
RUN go mod download

# Build the application
RUN go build -o profile-service ./cmd/server

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/profile-service/profile-service .
COPY --from=builder /app/profile-service/configs/config.yaml .
COPY --from=builder /app/profile-service/migrations /app/migrations
EXPOSE 4050
CMD ["./profile-service"]
