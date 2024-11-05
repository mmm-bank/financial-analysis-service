FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /build/transaction_consumer/mongo
RUN go build -o consumer .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/transaction_consumer/mongo/consumer .

ENTRYPOINT ["/app/consumer"]
