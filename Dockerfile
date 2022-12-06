FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/mhrynenko/jwt-service
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/github.com/mhrynenko/jwt-service /go/src/github.com/mhrynenko/jwt-service


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/github.com/mhrynenko/jwt-service /usr/local/bin/github.com/mhrynenko/jwt-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["github.com/mhrynenko/jwt-service"]
