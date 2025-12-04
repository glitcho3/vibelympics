#!/bin/sh
# Accept connections on 9001 and forward to sanitizer
echo "[$(date)]: entry" >&2
socat TCP-LISTEN:9001,reuseaddr,fork TCP:gawk:9002


