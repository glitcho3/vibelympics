#!/bin/sh
set -eu

if [ -d /output ]; then
    chown 65532:65532 /output 2>/dev/null || true
    chmod 755 /output 2>/dev/null || true
fi

echo "[$(date)] process: start" >&2
pandoc -f markdown -t html -o /output/index.html --include-in-header header.html
#pandoc -f markdown -t pdf -o /output/document.pdf # pritables flashcars, need some pdf render...
echo "[$(date)] process: end" >&2

