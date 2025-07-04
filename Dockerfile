# Future Docker image configuration
# Multi-stage build for Go application

# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o easypars cmd/easypars/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/easypars .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/frontend ./frontend

# Expose port
EXPOSE 8080

# Run the application
CMD ["./easypars"]