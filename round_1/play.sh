#!/bin/bash
set -eu

MOVIES_FILE="./movies.csv"
LOG="./game.log"
GUESS_TIME=3

#round=1

# leemos todas las líneas en un array
mapfile -t MOVIES < "$MOVIES_FILE"

shuffle() {
    MOVIES=( $(printf "%s\n" "${MOVIES[@]}" | shuf) )
}

shuffle

#i=0

while IFS=',' read -r name e1 e2 e3; do
    EMOJIS="$e1 $e2 $e3"
    echo "$EMOJIS" | nc -q 1 localhost 9001
    now=$(date '+%H:%M:%S')
    echo "$now round sent ($name): $EMOJIS"
    sleep $GUESS_TIME
done < "$MOVIES_FILE"
