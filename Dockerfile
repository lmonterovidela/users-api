FROM golang
ARG app_env

ENV ENV $app_env
COPY . /go/src/github.com/users-api
WORKDIR /go/src/github.com/users-api

RUN make build

CMD ./service

EXPOSE 8080
