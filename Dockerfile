FROM golang:1.24.5-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/main .
COPY internal/db/migrations ./internal/db/migrations
CMD ["./main"]
