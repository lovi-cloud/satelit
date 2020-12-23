FROM golang:1.15

MAINTAINER whywaita <https://github.com/whywaita>

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go get -u google.golang.org/grpc \
    && go get -u github.com/golang/protobuf/protoc-gen-go
RUN apt update -y \
    && apt install -y open-iscsi protobuf-compiler

WORKDIR /go/src/github.com/lovi-cloud/satelit

COPY ./satelit ./satelit

CMD /go/src/github.com/lovi-cloud/satelit/satelit
