FROM docker-registry.janusplatform.io/janus/golang-builder:latest@sha256:86c580812cc57a65080d01755029535d8291b451918180feb84fbe10e0b02081 AS builder

LABEL team="Team Rex"

ARG JANUS_STACK_NAME
ARG JANUS_STACK_VERSION
ARG JANUS_ENVIRONMENT
ARG GITHUB_TOKEN

ENV GROUP=apps
ENV USER=appuser

RUN apk update && apk add -q make curl curl-dev bash musl-dev

RUN ./update-git-config.sh

WORKDIR /app

COPY . .

RUN go mod download


RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o bin/service .

FROM alpine@sha256:124c7d2707904eea7431fffe91522a01e5a861a624ee31d03372cc1d138a3126

COPY --from=builder /app/bin/service /service

COPY . .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# required to run scratch container as non-privileged user
COPY --from=builder /etc/passwd /etc/passwd

USER appuser

EXPOSE 8080
EXPOSE 8443

CMD ["./service"]
