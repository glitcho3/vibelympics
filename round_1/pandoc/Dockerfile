# Stage 1: Build Pandoc
FROM cgr.dev/chainguard/wolfi-base as pandoc-builder

RUN apk add --no-cache libarchive-tools gzip libstdc++ curl ca-certificates rust

ADD https://github.com/jgm/pandoc/releases/download/3.8.2.1/pandoc-3.8.2.1-linux-amd64.tar.gz /tmp/pandoc.tar.gz
RUN mkdir -p /opt/pandoc \
    && bsdtar -xf /tmp/pandoc.tar.gz -C /opt/pandoc --strip-components=1

# Stage 2: Runtime
FROM cgr.dev/chainguard/wolfi-base as runtime
RUN apk add --no-cache socat 
WORKDIR /svc

COPY processor.sh .
COPY ingress.sh .
COPY header.html .

RUN chmod +x /svc/processor.sh \
    && chmod +x /svc/ingress.sh

#RUN chown 65532:65532 /output -R && \
#    chmod -R 755 /output 

#USER 65532

COPY --from=pandoc-builder /opt/pandoc/bin/pandoc /usr/local/bin/pandoc

EXPOSE 9003

ENTRYPOINT ["/svc/ingress.sh"]

