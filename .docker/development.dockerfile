FROM golang:1.11

RUN go get -v github.com/markbates/refresh

WORKDIR /go/src/github.com/essajiwa/teratur
