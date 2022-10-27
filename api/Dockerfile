FROM golang:1.16 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY chemical ./chemical
COPY db ./db
COPY server ./server

COPY labs.csv /mnt
COPY projects_041022.csv /mnt/projects.csv

RUN go build -o coshh

EXPOSE 8080

CMD ["/app/coshh"]
