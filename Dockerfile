FROM golang:1.14.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main_test.go Makefile ./
