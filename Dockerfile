FROM golang:1.21-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY config/envs/docker.env .env
COPY config/envs/docker.test.env test.env

RUN rm -r -f config/envs

RUN go build -o main cmd/server/main.go

CMD ["/app/main"]