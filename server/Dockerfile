# Build stage
FROM golang:1.24.4 as builder

WORKDIR /app

# Copy module files separately for better layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

# Cross-compile to amd64 for Debian runtime
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Runtime stage
FROM debian:bullseye-slim

# Install certs for HTTPS connections
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
