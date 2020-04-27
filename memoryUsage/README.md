# Data Type and Memory Usage

| 類型 | 大小 |
| - | - |
| bool | 1 byte |
| intN, uintN, floatN, complexN | N/8 byte |
| int, uint, uintptr | 1 platform |
| *T | 1 platform |
| sring | 2 platform (data, len) |
| []T | 3 platform (data, len,cap) |
| map | 1 platform |
| func | 1 platform |
| chan | 1 platform |
| interface | 2 platform (type, value) |

*platform: this type is usually 32 bits wide on 32-bit systems and 64 bits wide on 64-bit systems.

<br/>

# unsafe.Sizeof()

The size does not include any memory possibly referenced by x.

For instance, if x is a slice, Sizeof returns the size of the slice descriptor, not the size of the memory referenced by the slice.

<br/>

# Reference

1. [Go Example: Memory and Sizeof](https://dlintw.github.io/gobyexample/public/memory-and-sizeof.html)
2. [Golang 聖經: unsafe.Sizeof, Alignof 和 Offsetof](https://wizardforcel.gitbooks.io/gopl-zh/ch13/ch13-01.html)