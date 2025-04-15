FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build arguments for versioning
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

# Build the application with version information
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X github.com/chrisw-dev/golang-mock-oauth2-server/internal/version.Version=${VERSION} \
              -X github.com/chrisw-dev/golang-mock-oauth2-server/internal/version.Commit=${COMMIT} \
              -X github.com/chrisw-dev/golang-mock-oauth2-server/internal/version.BuildDate=${BUILD_DATE}" \
    -o mock-oauth2-server ./cmd/server

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