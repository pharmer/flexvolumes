FROM alpine

RUN set -x \
  && apk add --update --no-cache ca-certificates tzdata

COPY digitalocean /digitalocean
COPY digitalocean.sh /digitalocean.sh

CMD ["/digitalocean.sh"]