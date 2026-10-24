#!/bin/sh
# Additional arguments override defaults or select packages, e.g. ./as ./xmath.
set -eu

repo_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
cd "$repo_root"

if [ "$#" -eq 0 ]; then
    set -- ./...
fi

go test -run '^$' -bench . -benchmem -benchtime=200ms -count=3 "$@"
