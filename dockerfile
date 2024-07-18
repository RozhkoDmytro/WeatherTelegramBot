# Use the official Go image as the base image
FROM golang:1.22-alpine

# Install git
RUN apk update && apk add --no-cache git


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the .netrc file and set permissions
COPY .netrc /root/.netrc
RUN chmod 600 /root/.netrc

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
