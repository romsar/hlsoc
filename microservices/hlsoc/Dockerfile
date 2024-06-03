FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk update && apk add build-base

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ../.. .

RUN touch .env

RUN go build -o ./bin/serve ./cmd/serve

FROM alpine

WORKDIR /app

COPY --from=builder /app/bin/serve .
COPY --from=builder /app/*.env .

ENTRYPOINT ["./serve"]

EXPOSE 9090