FROM cgr.dev/chainguard/wolfi-base

RUN apk add --no-cache gawk socat

WORKDIR /svc

COPY sanitize.awk .
COPY emoji_db.awk .
COPY filter.awk .
COPY processor.sh .
COPY ingress.sh .

RUN chmod +x /svc/processor.sh \
    && chmod +x /svc/ingress.sh

EXPOSE 9002

ENTRYPOINT ["/svc/ingress.sh"]

