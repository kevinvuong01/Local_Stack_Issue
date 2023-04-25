FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod download


RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o bin/service .

FROM alpine@sha256:124c7d2707904eea7431fffe91522a01e5a861a624ee31d03372cc1d138a3126

COPY --from=builder /app/bin/service /service

COPY prime_api_seeds.yml prime_api_seeds.yml

EXPOSE 8080
EXPOSE 8443

CMD ["./service"]
