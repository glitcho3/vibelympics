FROM cgr.dev/chainguard/wolfi-base
RUN apk add --no-cache socat

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 9001
USER 65532

ENTRYPOINT ["/entrypoint.sh"]

