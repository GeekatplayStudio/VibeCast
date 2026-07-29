# Multi-stage Dockerfile for AgenticSFU Server & UI
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sfu-server ./cmd/sfu-server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/sfu-server .
COPY config.yaml .

EXPOSE 7880 8080 50000-60000/udp
CMD ["./sfu-server"]
