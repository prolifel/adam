FROM alpine:latest

RUN apk add --no-cache curl

COPY /cron/crontab /etc/crontabs/root

CMD ["crond", "-f"]