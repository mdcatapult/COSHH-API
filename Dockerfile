FROM golang:1.16 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY cmd ./cmd
COPY internal/chemical ./internal/chemical
COPY internal/db ./internal/db
COPY internal/server ./internal/server

COPY assets/labs.csv /mnt
COPY assets/projects_041022.csv /mnt/projects.csv

RUN go build -o coshh ./cmd/main.go

EXPOSE 8080

CMD ["/app/coshh"]
