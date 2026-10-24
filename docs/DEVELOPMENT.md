# Development

Run checks from any directory using the repository scripts:

```sh
sh scripts/check.sh
sh scripts/check.sh -race -covermode=atomic -coverprofile=/tmp/golibs.cover
go tool cover -func=/tmp/golibs.cover
go tool cover -html=/tmp/golibs.cover
```

`check.sh` checks formatting without changing files, then runs vet, tests, and
the build. Additional arguments are passed to `go test`. Tests use local HTTP
and SSH servers and need permission to bind loopback sockets. The optional
external RSS test requires `GOLIBS_NETWORK_TESTS=1` and is skipped in short mode.

CI runs Go 1.26 and stable Go on Linux, plus stable Go on macOS. The stable Linux
job includes race detection, a coverage artifact, and reachable-vulnerability
checks. To run the dependency checks locally:

```sh
go mod verify
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

To enable the optional pre-commit check, from the repository root:

```sh
ln -s ../../pre-commit .git/hooks/pre-commit
```

The hook runs the same check script and does not alter or stage files.

## Fuzzing

Run fuzz targets independently of benchmarks:

```sh
go test ./deque -run '^$' -fuzz '^FuzzDequeAgainstSlice$' -fuzztime=30s
go test ./structs -run '^$' -fuzz '^FuzzCombinedTagMetadata$' -fuzztime=30s
go test ./table -run '^$' -fuzz '^FuzzParseCSVInput$' -fuzztime=30s
```

The deque target compares operations against a bounded slice model. Struct tag
fuzzing checks termination and rule ordering; CSV fuzzing checks input handling
and preservation of the destination on errors.

## Benchmarks

```sh
sh scripts/bench.sh
sh scripts/bench.sh ./as ./xmath
sh scripts/bench.sh -benchtime=1s -count=5 ./set ./retry ./xio
```

The script defaults to three 200 ms samples and reports time, bytes, and
allocations per operation. Later flags override the defaults. Compare the same
workloads on the same Go version and machine, with other builds and tests
stopped. Save outputs outside the checkout and use repeated samples before
making latency claims. See [benchmark notes](BENCHMARKS.md).
