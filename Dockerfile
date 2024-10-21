# Use the official Golang base image
FROM golang:1.23.2 AS builder

# Set environment variables for Go
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy the CSV file into the image
COPY BoulderTrailHeads.csv .

# Build the server
RUN go build -o trail-finder ./main.go


# Create a small image for production
FROM alpine:3.18

# Set the working directory
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/trail-finder .
COPY --from=builder /app/BoulderTrailHeads.csv .

# Expose port 8080 for the application
EXPOSE 8080

# Run the application
CMD ["./trail-finder", "--server"]
