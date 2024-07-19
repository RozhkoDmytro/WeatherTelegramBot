# Build stage: Build the Go application
FROM golang:1.22-alpine AS builder

# Install git
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the .netrc file and set permissions
COPY .netrc /root/.netrc
RUN chmod 600 /root/.netrc

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Deploy stage: Create a minimal runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Copy the .env file from the build context to the runtime image
COPY .env .

# Command to run the executable
CMD ["./main"]
