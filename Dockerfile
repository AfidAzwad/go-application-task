# Step 1: Build the Go application
FROM golang:1.21-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o /app/main ./cmd/main.go

# Step 2: Run the Go application
FROM alpine:latest

# Install any necessary dependencies (like certs and libraries)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
