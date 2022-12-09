FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/distributed_lab/Auth
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/gitlab.com/distributed_lab/Auth /go/src/gitlab.com/distributed_lab/Auth


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/gitlab.com/distributed_lab/Auth /usr/local/bin/gitlab.com/distributed_lab/Auth
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["gitlab.com/distributed_lab/Auth"]
