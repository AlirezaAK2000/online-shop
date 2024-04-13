# Start from a smaller base image like Alpine
FROM golang:alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GONOPROXY="*"

# Set the working directory inside the container
WORKDIR /online-shop

# Install Git
RUN apk add --no-cache git

# Copy go.mod and go.sum files to cache dependencies downloading
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Build the Go online-shop
RUN go build -o online-shop .

# Start a new stage from a small base image
FROM scratch

# Copy the pre-built binary from the builder stage to the root directory inside the container
COPY --from=builder /online-shop/online-shop /

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/online-shop"]
