# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies if needed
RUN apk add --no-cache git

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# We'll build a single binary for cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/finapp-api ./cmd/api/main.go

# Final Stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/finapp-api .
# Copy migrations (important for startup)
COPY --from=builder /app/migrations ./migrations
# Copy config if needed, though we'll use env vars
COPY --from=builder /app/config.yaml .

# Add entrypoint script
COPY --from=builder /app/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./finapp-api"]
