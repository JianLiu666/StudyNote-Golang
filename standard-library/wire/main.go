package main

import (
	"fmt"
	"os"
)

func main() {
	event, err := InitializeEvent()
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}

	event.Start()
}
