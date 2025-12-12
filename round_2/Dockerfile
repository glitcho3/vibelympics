# Builder
FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /work

# Keep consistent cache path
ENV GOMODCACHE=/go/pkg/mod

# Copy go.mod and go.sum first for maximum cache
COPY go.mod go.sum ./

# Download deps with cache
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the source only after deps resolved
COPY . .

# Ensure Go build cache directory exists
RUN mkdir -p /root/.cache/go-build

# Build binaries using build cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o auditor-runner ./cmd/runner

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o runner-aggregator ./cmd/aggregator

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o auditor-provenor ./cmd/provenor

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o helm-auditor ./cmd/auditor

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o auditor-reporter ./cmd/reporter

# Runtime image
FROM cgr.dev/chainguard/static:latest

WORKDIR /

COPY --from=builder /work/auditor-runner .
COPY --from=builder /work/runner-aggregator .
COPY --from=builder /work/auditor-provenor .
COPY --from=builder /work/helm-auditor .
COPY --from=builder /work/auditor-reporter .

ENTRYPOINT ["/helm-auditor"]

