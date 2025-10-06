# Build stage
FROM golang:1.25.1-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY src/go.mod src/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o product-api .

# Run stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/product-api .

# Expose port
EXPOSE 8080

# Set environment variable
ENV PORT=8080

# Run the binary
CMD ["./product-api"]
