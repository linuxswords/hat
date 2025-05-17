# Use the official Golang image to build the application
FROM golang:1.24.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the application
RUN go build -o hat cmd/hat/main.go

# Use a minimal base image to run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/hat .

# Expose the application port
EXPOSE 8987

# Command to run the application
CMD ["./hat"]
