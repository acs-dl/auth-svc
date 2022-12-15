FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/distributed_lab/acs/auth
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/auth /go/src/gitlab.com/distributed_lab/acs/auth


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/auth /usr/local/bin/auth
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["auth"]
