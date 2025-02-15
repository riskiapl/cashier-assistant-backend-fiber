# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Add git and build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY .env .env

# Add CA certificates
RUN apk add --no-cache ca-certificates

# Expose port
EXPOSE 8000

# Run the application
CMD ["./main"]