package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

type event struct {
}

type item struct {
	kind string
}

var items = map[string]item{
	"a": {kind: "gopher"},
	"b": {kind: "rabbit"},
}

func doSlowThing() {
	time.Sleep(10 * time.Millisecond)
}

func fetch(name string) (item, error) {
	i, ok := items[name]
	doSlowThing()
	if !ok {
		return item{}, errors.New("item not found")
	}
	return i, nil
}

func match(pattern string) (names []string) {
	for name := range items {
		if ok, _ := filepath.Match(pattern, name); ok {
			names = append(names, name)
		}
	}
	return names
}

// Internal, Caller-side concurrency: channels
// keeping the channel local to the caller function make its usage much easier to see.
func query(pattern string) []item {
	names := match(pattern)

	c := make(chan item)
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			item, err := fetch(name)
			if err == nil {
				c <- item
			}
		}(name)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var items []item
	for i := range c {
		items = append(items, i)
	}
	return items
}

func main() {
	fmt.Println(query("[ab]*"))
}
