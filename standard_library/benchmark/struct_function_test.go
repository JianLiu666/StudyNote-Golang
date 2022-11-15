package pprof

import "testing"

type Object struct {
	Placeholder1 uint64
	Placeholder2 uint64
	Placeholder3 uint64
	Placeholder4 uint64
}

func (o Object) Get() *Object {
	return &o
}

func (o *Object) GetItself() *Object {
	return o
}

func BenchmarkDereference(b *testing.B) {
	o := &Object{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o = o.Get()
	}
	b.StopTimer()

	b.ReportAllocs()
}

func BenchmarkReference(b *testing.B) {
	o := &Object{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o = o.GetItself()
	}
	b.StopTimer()

	b.ReportAllocs()
}
