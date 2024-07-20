FROM golang:1.21.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN env GOOS=linux CGO_ENABLED=0 go build -o /command_message_service ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /command_message_service .


EXPOSE 3002

CMD ["./command_message_service"]
