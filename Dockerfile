FROM golang
ENV GO111MODULE=on

COPY . /go/src/github.com/users-api
WORKDIR /go/src/github.com/users-api

RUN go build cmd/*.go

CMD ./service

EXPOSE 8080