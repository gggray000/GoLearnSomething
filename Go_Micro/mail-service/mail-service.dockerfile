# base go image
FROM golang:1.24-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=direct
ENV GOSUMDB=off
RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

RUN chmod +x /app/mailerApp

# build a tiny docker imageß
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailerApp /app
COPY --from=builder /app/templates /templates

CMD ["/app/mailerApp"]