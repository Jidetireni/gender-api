# Build stage
FROM golang:1.25-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache ca-certificates tzdata git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Install goose for database migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the Go app with CGO disabled for a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd

# Final stage (Lightweight Alpine image)
FROM alpine:latest

# Set the current working directory
WORKDIR /app

# Copy the Pre-built binary file and CA certificates from the previous stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/server .
COPY --from=builder /app/internals/sql/migrations ./internals/sql/migrations

# Accept DB_URL from build environment
ARG DB_URL

# Set environment variables for goose
ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=${DB_URL}
ENV GOOSE_MIGRATION_DIR=/app/internals/sql/migrations

# Run migrations before starting the application
CMD ["sh", "-c", "goose -v up && ./server"]
