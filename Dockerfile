FROM golang:1.21-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

COPY cmd cmd
COPY docs docs
COPY internal internal
COPY migrations migrations
COPY pkg pkg
COPY test test

COPY config/config.go config/config.go
COPY config/envs/docker.env .env
COPY config/envs/docker.test.env test.env

COPY certs certs

RUN go build -o main cmd/server/main.go

EXPOSE 3000

CMD ["/app/main"]