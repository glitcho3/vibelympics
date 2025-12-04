#!/bin/sh
set -eu
echo "[$(date)] process: start" >&2

tee /tmp/debug.log \
  | gawk -f /svc/sanitize.awk -f /svc/filter.awk \
  | socat - TCP:pandoc:9003

echo "[$(date)] process: end" >&2

##!/bin/sh
#set -eu
#echo "[$(date)] process: start" >&2
#tee /tmp/debug.log | gawk -f /svc/filter.awk | socat - TCP:pandoc:9003
#echo "[$(date)] process: end" >&2


