# Dynamic Select Statement

施工中 ...

## Benchmark

```
go test -bench=. -run=none -benchtime=1000x -cpu=1 -benchmem -cpuprofile cpu.profile -memprofile mem.profile

BenchmarkReflectSelect      1000         431479472 ns/op       385178044 B/op    5823505 allocs/op
BenchmarkGoSelect           1000          28954545 ns/op           17262 B/op        404 allocs/op
BenchmarkSanity             1000             32247 ns/op               0 B/op          0 allocs/op
```