FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mock-oauth2-server ./cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/mock-oauth2-server .

# Expose the port that the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./mock-oauth2-server"]