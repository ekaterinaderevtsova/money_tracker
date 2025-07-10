FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .env
COPY internal/db/migrations ./internal/db/migrations
CMD ["./main"]
