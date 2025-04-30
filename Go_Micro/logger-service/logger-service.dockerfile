## base go image
#FROM golang:1.18-alpine AS builder
#
#RUN mkdir /app
#
#COPY . /app
#
#WORKDIR /app
#
#RUN apk add --no-cache git
#
#ENV GOPROXY=direct
#ENV GOSUMDB=off
#RUN CGO_ENABLED=0 go build -o loggerServiceApp ./cmd/api
#
#RUN chmod +x /app/loggerServiceApp
#
## build a tiny docker imageß
FROM alpine:latest

RUN mkdir /app

COPY loggerServiceApp /app

CMD ["/app/loggerServiceApp"]