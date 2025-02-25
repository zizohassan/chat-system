# Use the official Golang image as a build stage
FROM golang:1.22-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install reflex for live reload
RUN go install github.com/cespare/reflex@latest

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main cmd/server/main.go

# Use the official Alpine image for a lean production container
FROM alpine:latest

# Install Go in the final stage
RUN apk add --no-cache go

# Set the Current Working Directory inside the container
WORKDIR /chat-system/

# Copy the Pre-built binary file from the build stage
COPY --from=build /app/main .

# Copy the reflex binary
COPY --from=build /go/bin/reflex /usr/local/bin/reflex

# Copy reflex.conf for live reload configuration
COPY reflex.conf /chat-system/reflex.conf

# Copy the Go modules
COPY go.mod go.sum /chat-system/

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["reflex", "-c", "/chat-system/reflex.conf"]
