#!/bin/sh
set -eu

repo_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
cd "$repo_root"

unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
    printf 'Run gofmt on these files:\n%s\n' "$unformatted" >&2
    exit 1
fi

go vet ./...
go test "$@" ./...
go build ./...
