# Dockerfile
FROM golang:1.23.1-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main.go .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Expose port
EXPOSE 8080

# Run the binary
CMD ["go", "run", "main.go"]