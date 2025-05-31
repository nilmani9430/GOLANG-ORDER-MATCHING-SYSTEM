# Use the official Golang image as a build environment
FROM golang:1.24.3 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files (if exists) and download dependencies
COPY go.mod ./
COPY go.sum* ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .
# COPY ./.env  /.env 

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api main.go

# Use a lightweight base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=build /app .

# Expose the port your API listens on
EXPOSE 8080

# Command to run the application
CMD ["./api"]