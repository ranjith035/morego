# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make build-base

WORKDIR /app

# Copy Go workspace and module configurations
COPY go.work ./
COPY cmd/go.mod ./cmd/
COPY core/go.mod ./core/
COPY drivers/go.mod ./drivers/
COPY pkg/go.mod ./pkg/

# Fetch dependencies (if any)
RUN go work sync

# Copy source code
COPY cmd/ ./cmd/
COPY core/ ./core/
COPY drivers/ ./drivers/
COPY pkg/ ./pkg/

# Build the executable
RUN go build -ldflags="-w -s" -o bin/mobile ./cmd/main.go

# Stage 2: Final minimal runtime image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates libc6-compat adb

WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/bin/mobile /usr/local/bin/mobile

# Expose gRPC port (default port 50051)
EXPOSE 50051

# Set the binary as entrypoint
ENTRYPOINT ["/usr/local/bin/mobile"]
