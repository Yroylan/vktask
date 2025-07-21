FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go


FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app/server .
COPY .env .env

EXPOSE 8080

CMD ["./server"]