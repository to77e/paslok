FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y xclip

RUN touch paslok.db

RUN go build -o paslok ./cmd/paslok/main.go
