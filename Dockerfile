FROM golang:1.16-alpine

RUN apk update && apk add make gcc bash curl musl-dev

WORKDIR /go/src/github.com/Gobonoid/form
ADD . /go/src/github.com/Gobonoid/form

CMD make lint && make test
