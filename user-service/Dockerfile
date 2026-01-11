# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o user-service ./cmd/server

# Final minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
COPY configs/config.yaml .
COPY migrations /app/migrations
EXPOSE 4050
CMD ["./user-service"]