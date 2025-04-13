FROM golang:1.24 AS builder

WORKDIR /app

# Copy root go.mod/sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project context
COPY . .

# Build the appointment-service binary
# Ensure the path matches your project structure
RUN CGO_ENABLED=0 GOOS=linux go build -o appointment-service ./services/appointment-service/cmd/

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/appointment-service .

# Copy the service-specific .env file from the builder stage
COPY --from=builder /app/services/appointment-service/.env .

# WARNING: Embedding .env variables directly into the image is generally not recommended.
ENV $(cat .env | xargs)

# Command to run the service
CMD ["./appointment-service"] 