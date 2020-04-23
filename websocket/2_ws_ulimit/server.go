package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"syscall"

	"github.com/gorilla/websocket"
)

var connCount int64

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgradge connection
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
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		log.Printf("msg: %s", string(msg))
	}
}

// 每個 Process 在執行時系統都不會無限制的允許單個 Process 不斷的消耗資源, 因此都會設置資源限制。
// Linux 系統中使用 resource limit 來表示, 每個 Process 都可以設置不同的資源限制, 當前 Process
// 與底下的 Sub process 都會遵守此限制, 且其他的 Process 不會受到影響。
func main() {
	// Increase resources limitations
	// 描述資源軟體硬限制(resource limit)的結構體
	var rLimit syscall.Rlimit
	// syscall.RLIMIT_NOFILE 一個 Process 能打開的最大文件數, 預設是1024
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	// Soft Limit, 是指 Kernel 所能支持的資源上限
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

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
