# Start from a base Go image
FROM golang:1.23.2-bookworm

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the web server
EXPOSE 8080

# Command to run the application
CMD ["./main"]