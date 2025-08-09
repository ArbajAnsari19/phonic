# Orchestrator Service Dockerfile
# Main coordination service (STT → LLM → TTS)

FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build args
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the orchestrator service
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -a -installsuffix cgo \
    -o orchestrator ./services/orchestrator/

# =============================================================================
# Final stage
# =============================================================================

FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -S phonic && adduser -S phonic -G phonic

WORKDIR /home/phonic

# Copy binary
COPY --from=builder /app/orchestrator .

# Copy configs
COPY --from=builder /app/configs ./configs

# Change ownership
RUN chown -R phonic:phonic /home/phonic

# Switch to non-root user
USER phonic

# Expose port
EXPOSE 8084

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ./orchestrator --health-check || exit 1

# Start orchestrator
CMD ["./orchestrator"]
