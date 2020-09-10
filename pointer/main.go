package main

import (
	"fmt"
	"time"
)

type Object struct {
	UUID        string
	Description string
	NumberValue int
}

func NewObjectByPointer() *Object {
	return &Object{}
}

func (o *Object) SetUUIDByPointer(uuid string) {
	o.UUID = uuid
}

func (o *Object) SetDescriptionByPointer(text string) {
	o.Description = text
}
func (o *Object) SetNumberValueByPointer(number int) {
	o.NumberValue = number
}

func NewObjectByValue() Object {
	return Object{}
}

func (o Object) SetUUIDByValue(uuid string) Object {
	o.UUID = uuid
	return o
}

func (o Object) SetDescriptionByValue(text string) Object {
	o.Description = text
	return o
}
func (o Object) SetNumberValueByValue(number int) Object {
	o.NumberValue = number
	return o
}

func main() {
	times := 10000000
	benchmarkWithPointer(times)
	benchmarkWithValue(times)
}

func benchmarkWithPointer(times int) {
	obj := NewObjectByPointer()

	start := time.Now()
	for i := 0; i < times; i++ {
		obj.SetUUIDByPointer(fmt.Sprintf("random-generated-uid-%d", i))
		obj.SetDescriptionByPointer(fmt.Sprintf("%v", time.Now().UTC().Format("2006-01-02T15:04:05.999999-07:00")))
		obj.SetNumberValueByPointer(i)
	}
	elapsed := time.Since(start)

	fmt.Printf("WithPointer: %v\n", elapsed)
}

func benchmarkWithValue(times int) {
	obj := NewObjectByValue()

	start := time.Now()
	for i := 0; i < times; i++ {
		obj = obj.SetUUIDByValue(fmt.Sprintf("random-generated-uid-%d", i))
		obj = obj.SetDescriptionByValue(fmt.Sprintf("%v", time.Now().UTC().Format("2006-01-02T15:04:05.999999-07:00")))
		obj = obj.SetNumberValueByValue(i)
	}
	elapsed := time.Since(start)

	fmt.Printf("WithValue:   %v\n", elapsed)
}
