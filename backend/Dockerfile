# builder
FROM golang:1.23.0-alpine3.20 AS builder
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /app

COPY go.mod go.sum .
COPY . ./

RUN go mod tidy
RUN env CGO_ENABLED=1 go build  .

# container
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kube-dash-backend ./

EXPOSE 5000
RUN mkdir /config

CMD ["./kube-dash-backend", "-kubeconfig", "/config/kube-config.yml", "-secretkey", "/config/secret.key"]