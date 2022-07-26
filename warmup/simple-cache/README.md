# Design and implement a data structure for cache.

- `get(key)`
  - Get the value of the key if the key **exists** in the cache, otherwise return `-1`
- `put(key, value, weight)`
  - Set or insert the value, when the cache reaches its capacity, it should invalidate the least scored key.
  - The score is calculated as:
    - when current_time != last_access_time: `weight / ln(current_time - last_access_time + 1)`
    - else: `weight / -100`

Your data structure should be optimized for the computational complexity of `get(key)`
 - i.e. Average case for computational complexity of `get(key)` could be `O(1)`.

In your code, you can assume common data structure such as `array`, different type of `list`, `hash table` are available.

Please explain the computational complexity of `get(key)` and `put(...)` in **Big-O notation**.

# Benchmark Result

```
BenchmarkSampleSet_Size_1000-10     10000    63055.00 ns/op     96 B/op    1 allocs/op
BenchmarkSampleGet_Size_1000-10     10000       94.62 ns/op     23 B/op    1 allocs/op

BenchmarkSampleSet_Size_10000-10    10000      168.70 ns/op    180 B/op    2 allocs/op
BenchmarkSampleGet_Size_10000-10    10000      239.70 ns/op    184 B/op    3 allocs/op

BenchmarkSimpleSet_Size_1000-10     10000      214.40 ns/op    181 B/op    2 allocs/op
BenchmarkSimpleGet_Size_1000-10     10000      108.90 ns/op     23 B/op    1 allocs/op

BenchmarkSimpleSet_Size_10000-10    10000      168.50 ns/op    181 B/op    2 allocs/op
BenchmarkSimpleGet_Size_10000-10    10000      262.60 ns/op    184 B/op    3 allocs/op
```