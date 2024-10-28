# Use the official golang image as a base image
FROM golang:1.17-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o app .

# Start a new stage from scratch
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/app .

# Expose the port on which the application will listen
EXPOSE 8080

# Command to run the Go application
CMD ["./app"]
