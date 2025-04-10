# Start from the official Golang base image
FROM docker.io/golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o anichart .

FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/anichart .

COPY --from=builder /app/fonts ./fonts

EXPOSE 3000

CMD ["./anichart"]
