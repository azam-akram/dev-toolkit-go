# Use the official golang image as a base image

FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o http-server-go ./cmd

# Start a fresh runtime stage
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/http-server-go .

# Expose the port on which the application will listen
EXPOSE 8080

# Command to run the Go application
CMD ["./http-server-go"]
