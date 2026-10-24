# golibs

Small Go utility packages, imported individually. Requires Go 1.26 or newer.

[![Audit](https://github.com/SimonWaldherr/golibs/actions/workflows/audit.yml/badge.svg)](https://github.com/SimonWaldherr/golibs/actions/workflows/audit.yml)
[![Go Reference](https://pkg.go.dev/badge/simonwaldherr.de/go/golibs.svg)](https://pkg.go.dev/simonwaldherr.de/go/golibs)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![DOI](https://zenodo.org/badge/17567331.svg)](https://zenodo.org/badge/latestdoi/17567331)

```sh
go get simonwaldherr.de/go/golibs
```

```go
import "simonwaldherr.de/go/golibs/set"

languages := set.New("Go", "Rust", "Go")
shared := languages.Intersection(set.New("Go", "Java"))
fmt.Println(languages.Len(), shared.Has("Go")) // 2 true
```

See [compatibility notes](docs/COMPATIBILITY.md) for the Go requirement,
validation changes, certificate behavior, and SSH ownership rules.

## Packages

API documentation and tests live alongside the package sources. Executable
examples are in `example_test.go` files; larger examples are in [examples/](examples).

| Package | Purpose |
| --- | --- |
| [ansi](ansi) | Terminal colors and text styles |
| [arg](arg) | Command-line flags |
| [as](as) | Conversions between common types and string classification |
| [atomicfile](atomicfile) | Replace files through private temporary files |
| [bitmask](bitmask) | Read and change individual bits |
| [cache](cache) | In-memory caching |
| [cachedfile](cachedfile) | Cached file reads and writes |
| [channel](channel) | Fan-out channels |
| [csv](csv) | CSV loading and parsing |
| [deque](deque) | Generic double-ended queues |
| [env](env) | Typed environment variables and struct loading |
| [file](file) | Files, directories, and incremental reading |
| [foreach](foreach) | Iterate over collections and files |
| [gcurses](gcurses) | Terminal display updates |
| [gopath](gopath) | Executable and system paths |
| [gps](gps) | Geographic coordinates and distances |
| [graphics](graphics) | Image transformations and filters |
| [http](http) | HTTP request helpers |
| [log](log) | Application logging |
| [node](node) | Linked nodes |
| [pseudodb](pseudodb) | File-backed key/value storage |
| [re](re) | Repeated and scheduled calls |
| [regex](regex) | Regular expression helpers |
| [retry](retry) | Bounded retries with context cancellation and fixed delays |
| [rss](rss) | RSS and Atom feed parsing |
| [semaphore](semaphore) | Limit concurrent work with context-aware acquisition |
| [semver](semver) | Semantic version parsing and comparison |
| [set](set) | Generic sets and set operations |
| [specialdays](specialdays) | Calendar dates and holidays |
| [ssh](ssh) | SSH sessions with known-host verification |
| [ssl](ssl) | Self-signed server certificates and private keys |
| [stack](stack) | Stacks, queues, and rings |
| [structs](structs) | Struct reflection, validation, and modification |
| [table](table) | ASCII, Markdown, HTML, and CSV tables |
| [xio](xio) | Bounded stream reads and complete writes |
| [xmath](xmath) | Numeric aggregates and mathematical helpers |
| [xtime](xtime) | Time formatting with strftime-style patterns |
| [yoloDB](yoloDB) | In-memory document storage |

## Small utility packages

The six utilities below use only the standard library. Their examples cover
normal usage; API comments describe errors, ownership, and concurrency limits.

| Package | Main operations | Notes |
| --- | --- | --- |
| [atomicfile](atomicfile/example_test.go) | `WriteFile`, streaming `Write` | Syncs the file before rename; replaces a final symlink. Atomic visibility depends on OS rename semantics; the parent directory is not synced. |
| [deque](deque/example_test.go) | `PushFront`, `PushBack`, `PopFront`, `PopBack` | Zero value is ready. `Clone` copies storage; synchronize concurrent mutation externally. |
| [retry](retry/example_test.go) | `Do(ctx, attempts, delay, fn)` | Returns the last error unchanged. Cancellation interrupts delays; active callbacks must honor the context. |
| [semaphore](semaphore/example_test.go) | `New`, `Acquire`, `TryAcquire`, `Release` | Share one instance among workers. Each successful acquisition requires one release. |
| [set](set/example_test.go) | `Add`, `Has`, `Clone`, set algebra | Zero value is ready. Assignment shares initialized storage; `Clone` and set operations return independent storage. |
| [xio](xio/example_test.go) | `ReadLimit`, `WriteFull` | Reads consume one probe byte on overflow. Positive short writes are retried. Streams remain caller-owned. |

## Development

```sh
sh scripts/check.sh        # formatting, vet, tests, build
sh scripts/check.sh -race  # include the race detector
sh scripts/bench.sh ./as ./xmath
```

See [development instructions](docs/DEVELOPMENT.md) for coverage, fuzzing,
dependency checks, and hook setup, and [benchmark notes](docs/BENCHMARKS.md)
for workloads and allocation comparisons.
