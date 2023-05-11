FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/acs-dl/auth-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/auth-svc /go/src/github.com/acs-dl/auth-svc-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/auth /usr/local/bin/auth
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["auth"]
