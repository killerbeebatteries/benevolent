# based off https://firehydrant.com/blog/develop-a-go-app-with-docker-compose/
FROM golang:1.20 as base

FROM base as built

WORKDIR /go/app/
COPY *.go ./
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go build -o /tmp/irc-bot ./*.go

FROM alpine:latest as prod

WORKDIR /app
COPY --from=built /tmp/irc-bot /usr/bin/irc-bot
CMD ["/usr/bin/irc-bot"]
