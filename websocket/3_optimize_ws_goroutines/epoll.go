package main

import (
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
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &epoll{
		fd:          fd,
		connections: make(map[int]*websocket.Conn),
		lock:        &sync.RWMutex{},
	}, nil
}

func (e *epoll) Add(conn *websocket.Conn) error {
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	return nil
}

func websocketFD(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
