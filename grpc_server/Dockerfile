# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o grpcserver server.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/grpcserver .
COPY --from=builder /app/.env ./
COPY --from=builder /app/certs ./certs

EXPOSE 50051

# Run the server
CMD ["./grpcserver"]
