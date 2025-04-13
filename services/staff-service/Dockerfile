FROM golang:1.24 AS builder

WORKDIR /app

# Copy root go.mod/sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project context
# This is needed because the build path references the service directory
COPY . .

# Build the staff-service binary
# Ensure the path matches your project structure
RUN CGO_ENABLED=0 GOOS=linux go build -o staff-service ./services/staff-service/cmd/

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app # Changed to /app for consistency

# Copy the binary from the builder stage
COPY --from=builder /app/staff-service .

# Copy the service-specific .env file from the builder stage
COPY --from=builder /app/services/staff-service/.env .

# WARNING: Embedding .env variables directly into the image is generally not recommended.
# Consider loading at runtime instead.
ENV $(cat .env | xargs)

# Command to run the service
CMD ["./staff-service"] 