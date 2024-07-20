# MessagesService/Dockerfile
FROM golang:1.21.2 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./

# Build the Go app
RUN env GOOS=linux CGO_ENABLED=0 go build -o /command_message_service ./cmd/api

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /command_message_service .

# Expose port 8080 to the outside world
EXPOSE 3002

# Command to run the executable
CMD ["./command_message_service"]
