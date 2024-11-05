FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /build/transaction_consumer/postgres
RUN go build -o consumer .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/transaction_consumer/postgres/consumer .

ENTRYPOINT ["/app/consumer"]
