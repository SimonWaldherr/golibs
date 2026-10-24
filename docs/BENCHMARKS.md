# Benchmark notes

Run `sh scripts/bench.sh` for all packages, or pass package paths and Go test
flags. Benchmarks include allocation reporting and exclude ordinary tests.

## Workloads

- `as`: fixed string classifiers and accepted/rejected email and URL patterns.
- `xmath`: aggregates over 1,024 numeric values, with median as a copying control.
- `file`: delimiter reads of 4–16 KiB and writes of 2,048 small strings.
- `table`: rendering and CSV parsing of 100 rows.
- `structs`: repeated validation/modification of one six-field type, after its
  rule metadata and expressions have been cached.
- `set`: membership, overlapping algebra, and 1,024-element empty-operand cases.
- `retry`: first success, immediate exhaustion, and five failing attempts with
  four one-nanosecond delays. The delayed case measures timer overhead.
- `xio`: a bounded 4 KiB read and a complete write to `io.Discard`.
- `atomicfile`: replacement of 4 KiB, including file sync.
- `semaphore`: sequential and contended parallel acquire/release pairs.
- `deque`: a warm 1,024-element queue and growth from empty.

Discard writes measure helper overhead. File sync and contention benchmarks
depend strongly on the filesystem and scheduler.

## Local allocation comparisons

Measured with Go 1.27.3 on darwin/arm64 (Apple M2 Max), using three 200 ms
samples. Existing-package baselines use the implementations at `ea1378a` with
the new benchmark files added. Set and retry compare the initial utility
implementations with their final clone/timer optimizations.

| Operation | Allocations/op before → after |
| --- | ---: |
| Email validation, accepted or rejected | 64 → 0 |
| URL validation, accepted or rejected | 90 → 0 |
| ReadUntil, 16 KiB | 16,399 → 21 |
| HTML rendering, 100 rows | 1,504 → 504 |
| CSV parsing, 100 rows | 430 → 323 |
| ValidateAndModify, six fields | 576 → 3 |
| Union with an empty receiver, 1,024 keys | 22 → 6 |
| Difference with an empty operand, 1,024 keys | 22 → 6 |
| Retry, five delayed failing attempts | 12 → 3 |

Numeric reductions other than median, successful first-attempt retries, set
membership, warm deque operations, and semaphore acquire/release do not allocate in these
workloads. Builders remove repeated string copying from delimiter reads;
buffered line writes trade a 4 KiB buffer for fewer write calls.

System load varied during measurement, so these results emphasize allocations.
Re-run on an idle machine with the same Go version and inputs before drawing
latency conclusions. To focus on the latest optimizations:

```sh
go test -run '^$' -bench 'Benchmark(SetUnionEmpty|SetDifferenceEmpty|RetryDelayed)$' \
  -benchmem -benchtime=1s -count=5 ./set ./retry
```
