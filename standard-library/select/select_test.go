package dynamicselect

import (
	"reflect"
	"testing"
)

const numProduce = 1000
const numChannels = 100
const expectResult int64 = numProduce * numChannels * (numProduce - 1) / 2

func BenchmarkReflectSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkResult(b, benchmarkReflectSelect())
	}
}

func BenchmarkGoSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkResult(b, benchmarkGoSelect())
	}
}

func BenchmarkSanity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkResult(b, benchmarkSanity())
	}
}

// use reflection to select from a slice of channels
func benchmarkReflectSelect() int64 {
	channels := makeChannels(numChannels)
	var result int64 = 0

	cases := make([]reflect.SelectCase, numChannels)
	for i, ch := range channels {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	for finished := 0; finished < numChannels; {
		i, value, ok := reflect.Select(cases)
		if !ok {
			cases = append(cases[:i], cases[i+1:]...)
			finished++
		} else {
			result += value.Int()
		}
	}

	return result
}

// use an aggregate channel to select from a slice of channels
func benchmarkGoSelect() int64 {
	channels := makeChannels(numChannels)
	var result int64 = 0

	done := make(chan struct{})
	combinedChannel := make(chan int)
	for i := 0; i < numChannels; i++ {
		go func(c chan int) {
			for v := range c {
				combinedChannel <- v
			}
			done <- struct{}{}
		}(channels[i])
	}

	finished := 0
	for finished < numChannels {
		select {
		case i := <-combinedChannel:
			result += int64(i)
		case <-done:
			finished++
		}
	}

	close(combinedChannel)
	close(done)

	return result
}

func benchmarkSanity() int64 {
	var result int64 = 0
	for i := 0; i < numChannels; i++ {
		for j := 0; j < numProduce; j++ {
			result += int64(j)
		}
	}

	return result
}

func checkResult(b *testing.B, result int64) {
	if result != expectResult {
		b.Fatalf("Fail! Expected %v but got %v", expectResult, result)
	}
}

func produce(c chan int, n int) {
	for i := 0; i < n; i++ {
		c <- i
	}
	close(c)
}

func makeChannels(n int) []chan int {
	chans := make([]chan int, n)
	for i := 0; i < n; i++ {
		c := make(chan int)
		go produce(c, numProduce)
		chans[i] = c
	}

	return chans
}
