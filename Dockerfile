FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .env
CMD ["./main"]
