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