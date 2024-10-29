# Build stage
FROM golang:1.23-alpine AS builder

# Environment variables for static linking
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application as a statically linked binary
RUN go build -o app ./cmd/app

# Final stage - use minimal base image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install only necessary tools
RUN apk add --no-cache bash ca-certificates

# Copy the compiled binary and .env file from the build stage
COPY --from=builder /app/app .
COPY .env .env

# Run the app binary with .env loaded
CMD ["sh", "-c", "source .env && ./app"]
