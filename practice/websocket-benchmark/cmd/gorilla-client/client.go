package main

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"
	"websocketbenchmark/model"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type client struct {
	conn  *websocket.Conn    //
	times []map[string]int64 //
}

func newClient(u url.URL) *client {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Fatal("dial:", err)
	}

	return &client{
		conn:  conn,
		times: make([]map[string]int64, conf.Simulation.NumMessages),
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

			var payload model.Payload
			err = json.Unmarshal(message, &payload)
			if err != nil {
				logrus.Error("unmarshal:", err)
				return
			}

			c.times[payload.Count]["server_received"] = payload.Timestamp
			c.times[payload.Count]["client_received"] = time.Now().UnixMilli()

			if payload.Count >= int64(conf.Simulation.NumMessages)-1 {
				wg.Done()
				return
			}
		}
	}()

	for i := 0; i < conf.Simulation.NumMessages; i++ {
		c.times[i] = map[string]int64{
			"client_start": time.Now().UnixMilli(),
		}

		payload := model.Payload{
			Count:     int64(i),
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
