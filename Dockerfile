FROM golang:1.5.1
MAINTAINER Matthias Luedtke (matthiasluedtke)

RUN apt-get update
RUN apt-get install -y \
 imagemagick \
 pkg-config \
 libmagickwand-dev

COPY . /go/src/github.com/mat/gomagick/
WORKDIR /go/src/github.com/mat/gomagick
RUN go get ./...

EXPOSE 8080

ENV PORT 8080
ENTRYPOINT "magickserver"
