# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/idxstocks ./cmd/api/main.go

# Final stage
FROM alpine:latest

# Install essential packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/idxstocks .

# Copy migrations folder
COPY --from=builder /app/migrations ./migrations

# Expose the application port
EXPOSE 3000

# Set environment variables (Defaults)
ENV PORT=3000
ENV TZ=Asia/Jakarta

# Run the binary
CMD ["./idxstocks"]
