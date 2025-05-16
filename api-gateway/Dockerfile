# Build stage
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./cmd/main.go

# Final stage
FROM alpine:3.17

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set the timezone
ENV TZ=UTC

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api-gateway .

# Use the non-root user
USER appuser

# Expose port
EXPOSE 8000

# Command to run the executable
ENTRYPOINT ["./api-gateway"]