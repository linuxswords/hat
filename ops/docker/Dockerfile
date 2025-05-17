# Use the official Golang image to build the application
FROM golang:1.24.2-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . .

# Install build dependencies
# RUN apk add --no-cache gcc musl-dev sqlite-dev
RUN apt-get update && apt-get install -y gcc libsqlite3-dev

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o hat ./cmd/hat/

# Use a minimal base image to run the application
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Install required runtime libraries
RUN apt-get update && apt-get install -y libsqlite3-0 && rm -rf /var/lib/apt/lists/*
RUN mkdir -p internal/models/bootstrap
# Copy the binary from the builder stage
COPY --from=builder /app/hat .
COPY --from=builder /app/static/ static/
COPY --from=builder /app/templates/ templates/
COPY --from=builder /app/internal/models/boostrap/*sql internal/models/bootstrap/

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./hat"]
