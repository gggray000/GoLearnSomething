# base go image
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN apk add --no-cache git

#ENV GOPROXY=direct
#ENV HTTP_PROXY=
#ENV HTTPS_PROXY=
#ENV http_proxy=
#ENV https_proxy=

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# build a tiny docker imageß
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]