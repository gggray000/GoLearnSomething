FROM alpine:latest

RUN mkdir /app

COPY frontEndLinuxApp /app

CMD ["/app/frontEndLinuxApp"]