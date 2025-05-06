# Stage 1: Build the Go binary
FROM golang:1.23.5-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o canvas-service

# Stage 2: Minimal image
FROM alpine:latest

WORKDIR /root/

# Копируем собранный бинарник
COPY --from=builder /app/canvas-service .

# По умолчанию запускаем наше приложение
CMD ["./canvas-service"]
