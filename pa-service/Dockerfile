# Start from the official Golang base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o service

# Expose the desired port
EXPOSE 8000

# Set the command to run the binary
CMD ["./service"]