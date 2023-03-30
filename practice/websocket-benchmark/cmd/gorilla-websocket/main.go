package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"websocketbenchmark/internal/config"
	"websocketbenchmark/model"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "net/http/pprof"
)

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})

	viper.SetConfigFile("./conf.d/env.yaml")
	viper.AutomaticEnv()
}

func main() {
	conf := config.NewFromViper()

	http.HandleFunc("/echo", echo)

	addr := fmt.Sprintf("%s:%s", conf.Server.Addr, conf.Server.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}

	//upgrade the connection from a HTTP connection to a websocket connection
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("websocket upgrade:", err)
		return
	}
	defer c.Close()

	err = notify(c, 0)
	if err != nil {
		logrus.Error("notify:", err)
		return
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logrus.Error("read:", err)
			break
		}

		// logrus.Info("recv:", string(message))
		if mt != websocket.TextMessage {
			continue
		}

		var json_data model.Payload
		err = json.Unmarshal(message, &json_data)
		if err != nil {
			logrus.Error("json unmarshal:", err)
			return
		}

		err = notify(c, json_data.Count)
		if err != nil {
			logrus.Error("notify:", err)
			return
		}
	}
}

// Send a connected client an event JSON string
// @param ws - The client connection the outgoing message is for
// @param c  - The message count
//
// @return Error object containing a possible error that occured
func notify(ws *websocket.Conn, c int32) error {
	return ws.WriteMessage(websocket.TextMessage, getEvent(c))
}

// Creates a JSON string containing the message count and the current timestamp
// @param c - The message count
//
// @return A JSON string (byte array) containing the message count and the current timestamp
func getEvent(c int32) []byte {
	var event model.Payload
	event.Count = c
	event.Timestamp = time.Now().UnixMilli()

	b, err := json.Marshal(event)
	if err != nil {
		logrus.Error("json marshal failed:", err)
		return []byte{}
	}

	return b
}
