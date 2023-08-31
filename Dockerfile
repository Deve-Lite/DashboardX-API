FROM golang:1.21-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY .env.docker .env

RUN go build -o main cmd/server/main.go

CMD ["/app/main"]