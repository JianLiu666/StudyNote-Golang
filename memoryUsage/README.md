# unsafe.Sizeof()

The size does not include any memory possibly referenced by x.

For instance, if x is a slice, Sizeof returns the size of the slice descriptor, not the size of the memory referenced by the slice.

# Reference

1. [Go Example: Memory and Sizeof](https://dlintw.github.io/gobyexample/public/memory-and-sizeof.html)