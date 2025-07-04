# ==========================
# Build stage
# ==========================
FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

# Build a fully static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# ==========================
# Runtime stage
# ==========================
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates dnsutils
WORKDIR /app

# Copy the statically linked binary
COPY --from=builder /app/app .

# Copy CA certs so HTTPS works
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["./app"]

