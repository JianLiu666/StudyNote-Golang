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
BenchmarkSampleSet_Size_1000-10     20000     67198.00 ns/op     90 B/op    1 allocs/op
BenchmarkSampleGet_Size_1000-10     20000        86.99 ns/op     13 B/op    0 allocs/op

BenchmarkSampleSet_Size_10000-10    20000    352370.00 ns/op    133 B/op    2 allocs/op
BenchmarkSampleGet_Size_10000-10    20000       158.80 ns/op     94 B/op    1 allocs/op

BenchmarkSimpleSet_Size_1000-10     20000       224.90 ns/op    181 B/op    2 allocs/op
BenchmarkSimpleGet_Size_1000-10     20000        91.42 ns/op     13 B/op    0 allocs/op

BenchmarkSimpleSet_Size_10000-10    20000       188.80 ns/op    182 B/op    2 allocs/op
BenchmarkSimpleGet_Size_10000-10    20000       166.00 ns/op     94 B/op    1 allocs/op
```