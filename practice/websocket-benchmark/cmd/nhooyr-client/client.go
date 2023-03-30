package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"websocketbenchmark/model"

	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
)

type client struct {
	conn  *websocket.Conn                   //
	times [numMessages + 1]map[string]int64 //
}

func newClient(ctx context.Context, u string) *client {
	conn, _, err := websocket.Dial(ctx, u, &websocket.DialOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		logrus.Fatal("dial: ", err)
	}

	return &client{
		conn:  conn,
		times: [numMessages + 1]map[string]int64{{}},
	}
}

func (c *client) start(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		for {
			_, message, err := c.conn.Read(ctx)
			if err != nil {
				logrus.Error("read:", err)
				return
			}
			bar.Add(1)

			var payload model.Payload
			err = json.Unmarshal(message, &payload)
			if err != nil {
				logrus.Error("unmarshal: ", err)
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

		payload := model.Payload{
			Count:     int32(i),
			Timestamp: 0,
		}
		binary, err := json.Marshal(&payload)
		if err != nil {
			logrus.Error("marshal:", err)
			return
		}
		err = c.conn.Write(ctx, websocket.MessageText, binary)
		if err != nil {
			logrus.Error("write:", err)
			return
		}
	}
}

func (c *client) close() {
	c.conn.Close(websocket.StatusNormalClosure, "")
}
