package main

import (
	"log"
	"reflect"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
)

type epoll struct {
	fd          int
	connections map[int]*websocket.Conn
	lock        *sync.RWMutex
}

func MkEpoll() (*epoll, error) {
	// 過去的作法是 EpollCreate 建立一個實例時需要告訴 kernel 要監聽 size 個 fd, 當使用完畢後
	// 如果沒有調用 close(), 可能會將 process 可以使用的 fd 資源耗盡
	// Linux 2.6.8 版本後, kernel 可以動態分配大小, 因此不需要在特地關心這個 size 了
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]*websocket.Conn),
	}, nil
}

func (e *epoll) Add(conn *websocket.Conn) error {
	fd := websocketFD(conn)

	// 註冊目標 fd 到 epfd 中, 同時關聯內部 event 到 fd 上
	// events 參數是一個枚舉的集合, 可以用 `|` 來增加事件類型, 這裡帶入的兩個參數為
	// EPOLLIN: 表示關聯的 fd 可以進行讀操作
	// EPOLLHUP: 表示關聯的 fd 已掛起, epoll_wait 會一直等待這個事件, 所以沒必要設置這個屬性
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	e.connections[fd] = conn
	if len(e.connections)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.connections))
	}

	return nil
}

func (e *epoll) Remove(conn *websocket.Conn) error {
	fd := websocketFD(conn)

	// 從 epfd 中刪除已註冊的 fd, event 可以被忽略(nil)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	delete(e.connections, fd)
	if len(e.connections)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.connections))
	}

	return nil
}

func (e *epoll) Wait() ([]*websocket.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}

	e.lock.RLock()
	defer e.lock.RUnlock()

	var connections []*websocket.Conn
	for i := 0; i < n; i++ {
		conn := e.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}

	return connections, nil
}

func websocketFD(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
