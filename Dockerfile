# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

# Install build dependencies for CGO (required for sqlite3)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Build arguments for cross-compilation
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Enable CGO for sqlite3
ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with proper architecture
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o adam .

FROM alpine:latest
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

# Copy binary and migrations as root
COPY --from=builder /app/adam /app/adam
COPY --from=builder /app/migrations /app/migrations

# Create data directory
RUN mkdir -p /app/data

# Make binary executable
RUN chmod +x /app/adam

EXPOSE 8080
ENTRYPOINT ["/app/adam"]