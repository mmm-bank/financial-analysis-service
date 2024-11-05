FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o financial-analysis .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/financial-analysis .

ENTRYPOINT ["/app/financial-analysis"]
