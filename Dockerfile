# Multi-stage Dockerfile for Kortex Web Server

# Stage 1: Builder
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the web server binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o kortex-web ./cmd/web

# Stage 2: Runtime
FROM mcr.microsoft.com/playwright:v1.40.0-jammy

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/kortex-web .

# Copy .env.example as reference (users should mount their own .env)
COPY .env.example .

# Set default environment variables
ENV HEADLESS=true
ENV PORT=8080
ENV DB_PATH=/app/data/kortex.db

# Create data directory for database
RUN mkdir -p /app/data

# Expose the web server port
EXPOSE 8080

# Run the web server
CMD ["./kortex-web"]
