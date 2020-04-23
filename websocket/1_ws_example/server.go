package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var connCount int64

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	n := atomic.AddInt64(&connCount, 1)
	if n%100 == 0 {
		log.Printf("Total number of connection: %v", n)
	}
	defer func() {
		n := atomic.AddInt64(&connCount, -1)
		if n%100 == 0 {
			log.Printf("Total number of connection: %v", n)
		}
		conn.Close()
	}()

	// Read message from socket
	msgCount := uint64(0)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if msgCount%1000 == 0 {
			log.Printf("msg: %s", string(msg))
		}
		msgCount++
	}
}

func main() {
	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof failed: %v", err)
		}
	}()

	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
