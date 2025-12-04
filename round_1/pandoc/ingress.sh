#!/bin/sh
set -eu

exec socat TCP-LISTEN:9003,reuseaddr,fork \
     EXEC:"/svc/processor.sh",nofork
