package main

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type client struct {
	conn  *websocket.Conn                   //
	times [numMessages + 1]map[string]int64 //
}

func newClient(u url.URL) *client {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Fatal("dial:", err)
	}

	return &client{
		conn:  conn,
		times: [numMessages + 1]map[string]int64{{}},
	}
}

func (c *client) start(wg *sync.WaitGroup) {
	go func() {
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				logrus.Error("read:", err)
				return
			}
			bar.Add(1)

			var payload data
			err = json.Unmarshal(message, &payload)
			if err != nil {
				logrus.Error("unmarshal:", err)
				return
			}

			c.times[payload.Count]["server_received"] = payload.Timestamp
			c.times[payload.Count]["client_received"] = time.Now().UnixMilli()

			if payload.Count >= int32(numMessages) {
				wg.Done()
				return
			}
		}
	}()

	for i := 1; i <= numMessages; i++ {
		c.times[int32(i)] = map[string]int64{
			"client_start": time.Now().UnixMilli(),
		}

		payload := data{
			Count:     int32(i),
			Timestamp: 0,
		}
		binary, err := json.Marshal(&payload)
		if err != nil {
			logrus.Error("marshal:", err)
			return
		}
		err = c.conn.WriteMessage(websocket.TextMessage, binary)
		if err != nil {
			logrus.Error("write:", err)
			return
		}
	}
}

func (c *client) close() {
	c.conn.Close()
}
