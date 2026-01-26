# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY proto/ proto/
RUN go mod download

COPY . .
ENV PB_REPO https://github.com/protocolbuffers/protobuf
ENV PB_VERSION 30.2

RUN wget $PB_REPO/releases/download/v$PB_VERSION/protoc-$PB_VERSION-linux-x86_64.zip
RUN unzip protoc-$PB_VERSION-linux-x86_64.zip -d $HOME/protobuf
RUN mv $HOME/protobuf/bin/protoc /usr/local/bin

# Development stage
FROM builder AS development
COPY . .
RUN chmod +x ./regen.sh
RUN ./regen.sh

RUN go install github.com/air-verse/air@v1.62.0

CMD ["air", "-c", "air.toml"]

# Production build stage
FROM builder as production
COPY . .
RUN go build -o grpcserver ./cmd/server.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=production /app/grpcserver .
COPY --from=production /app/certs ./certs

EXPOSE 50051

# Run the server
CMD ["./grpcserver"]
