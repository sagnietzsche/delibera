#!/usr/bin/env bash
set -e

DATA_DIR="${1:-/tmp/delibera}"
NODES="${2:-1}"

cleanup() {
    echo "Shutting down..."
    kill 0
    rm -rf "$DATA_DIR"
}
trap cleanup EXIT INT TERM

go build -o delibera . 2>/dev/null

mkdir -p "$DATA_DIR"

echo "Starting node1 (leader)..."
./delibera -id node1 -haddr localhost:11000 -raddr localhost:12000 "$DATA_DIR/node1" &

if [ "$NODES" -ge 3 ]; then
    sleep 1
    echo "Starting node2..."
    ./delibera -id node2 -haddr localhost:11001 -raddr localhost:12001 -join localhost:11000 "$DATA_DIR/node2" &

    echo "Starting node3..."
    ./delibera -id node3 -haddr localhost:11002 -raddr localhost:12002 -join localhost:11000 "$DATA_DIR/node3" &
fi

echo ""
echo "Cluster up:"
echo "  node1  http://localhost:11000  raft: localhost:12000"
[ "$NODES" -ge 3 ] && echo "  node2  http://localhost:11001  raft: localhost:12001"
[ "$NODES" -ge 3 ] && echo "  node3  http://localhost:11002  raft: localhost:12002"
echo ""
echo "Press Ctrl+C to stop."

wait
