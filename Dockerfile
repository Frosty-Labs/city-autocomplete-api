# Use multi-stage build to keep the final image small
FROM golang:1.21-alpine AS builder

# Install git and sqlite3 dependencies for building
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

# Install sqlite package
RUN apk --no-cache add ca-certificates sqlite

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the binary, database file, and CSV file
COPY --from=builder /app/cities.db ./
COPY --from=builder /app/world-cities.csv ./

# Change ownership to appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose port (default is 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]