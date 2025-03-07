FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot-go .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/telegram-bot-go .
ENV BOT_TOKEN=""
CMD ["./telegram-bot-go"]