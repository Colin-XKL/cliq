# Multi-stage build for smaller final image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cliq-hub-backend ./cmd/server

# Final stage
FROM alpine:latest

# Install CA certificates for SSL/TLS connections
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/cliq-hub-backend .

# Copy templates directory if exists
COPY --from=builder /app/templates ./templates

# Change ownership to appuser
RUN chown -R appuser:appgroup /root/

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Environment variables with defaults
ENV GIN_MODE=release
ENV PORT=8080

# Run the binary
CMD ["./cliq-hub-backend"]