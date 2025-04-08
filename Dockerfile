# Use the official Golang image as the builder
FROM golang:1.18-buster AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/main .

# Command to run the application
CMD ["./main"]