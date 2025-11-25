# Multi-stage build for smaller image size

# Stage 1: Build the Go application
FROM golang:1.21-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    git \
    gcc \
    g++ \
    libc6-dev \
    libtesseract-dev \
    libleptonica-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o ocr-api .

# Stage 2: Create the runtime image
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    tesseract-ocr-tha \
    tesseract-ocr-eng \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Set timezone (optional)
ENV TZ=Asia/Bangkok

# Create app directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ocr-api .

# Create necessary directories
RUN mkdir -p /app/uploads

# Set environment variables
ENV SERVER_PORT=8077
ENV DATABASE_PATH=/app/data/db.sqlite
ENV UPLOAD_DIR=/app/uploads
ENV TESSERACT_LANG=tha+eng
ENV GIN_MODE=release

# Create volume mount point for database persistence
VOLUME ["/app/data"]

# Expose port
EXPOSE 8077

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./ocr-api"]
