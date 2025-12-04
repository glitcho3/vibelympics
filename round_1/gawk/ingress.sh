#!/bin/sh
set -eu

exec socat TCP-LISTEN:9002,reuseaddr,fork \
     EXEC:"/svc/processor.sh",nofork

