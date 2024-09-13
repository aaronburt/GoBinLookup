# Use a base image with Go installed
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy all files from the current directory into /app in the container
COPY . /app

# Install dependencies for CGO
RUN apk add --no-cache sqlite-dev gcc musl-dev

# Build the Go binary
RUN go build -o server .

# Use a minimal base image for the Go binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/server /app/server

# Copy additional files like .env from the builder stage
COPY --from=builder /app/.env /app/.env

# Copy the static folder and its contents from the builder stage
COPY --from=builder /app/static /app/static

# Make the Go binary executable
RUN chmod +x /app/server

# Execute the binary
CMD ["/app/server"]
