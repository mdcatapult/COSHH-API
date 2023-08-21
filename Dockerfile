FROM golang:1.17 AS builder

WORKDIR /app

COPY .env ./
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal/chemical ./internal/chemical
COPY internal/db ./internal/db
COPY internal/server ./internal/server

RUN go build -o coshh ./cmd/main.go

EXPOSE 8080

CMD ["/app/coshh"]
