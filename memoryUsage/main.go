package main

import (
	"fmt"
	"unsafe"
)

// Strings in Go are just headers containing a pointer and a length.
type StringHeader struct {
	Data uintptr
	Len  int
}

type T struct {
	B  uint8 // is a byte
	I  int   // it ts int32 on x86 PC or int64 on x64 PC
	P  *int  // it ts int32 on x86 PC or int64 on x64 PC
	S  string
	SS []string
}

func main() {
	fmt.Println("Size of type Int is", unsafe.Sizeof(int(0)))

	var text string = "memory usage"

	// So unsafe.Sizeof(string) will report the size of the above struct, which is independent
	// of the length of the `string` value(which is the value of the `Len` field)
	fmt.Println("Size of type String is", unsafe.Sizeof(text))

	// Go stores the UTF-8 encoded byte sequences of `string` values in memory. The builtin function
	// len() report the byte-length of a `string`, so basically the memory required to store a `string`
	// value in memory is:
	fmt.Println("Size of type String is", len(text)+int(unsafe.Sizeof(text)))

	fmt.Println("=")

	t := &T{}
	fmt.Printf("sizeof(uint8)    %2v, offset= %2v\n", unsafe.Sizeof(t.B), unsafe.Offsetof(t.B))
	fmt.Printf("sizeof(int)      %2v, offset= %2v\n", unsafe.Sizeof(t.I), unsafe.Offsetof(t.I))
	fmt.Printf("sizeof(*int)     %2v, offset= %2v\n", unsafe.Sizeof(t.P), unsafe.Offsetof(t.P))
	fmt.Printf("sizeof(string)   %2v, offset= %2v\n", unsafe.Sizeof(t.S), unsafe.Offsetof(t.S))
	fmt.Printf("sizeof([]string) %2v, offset= %2v\n", unsafe.Sizeof(t.SS), unsafe.Offsetof(t.SS))
	fmt.Printf("sizeof(T)        %2v\n", unsafe.Sizeof(t))
}
