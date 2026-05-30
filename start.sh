#!/usr/bin/env bash
set -e

DATA_DIR="${1:-/tmp/delibera}"
NODES="${2:-1}"

go build -o delibera . 2>/dev/null

./delibera start --data-dir "$DATA_DIR" --nodes "$NODES"
