package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"websocketbenchmark/config"
	"websocketbenchmark/model"
	"websocketbenchmark/util"

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
		logrus.Errorf("websocket upgrade: %v", err)
		return
	}
	defer c.Close()

	ch := make(chan []byte, 1000)
	go func() {
		defer close(ch)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logrus.Error("read: ", err)
				break
			}

			ch <- message
		}
	}()

	for payload := range ch {
		var recv model.Payload
		err = json.Unmarshal(payload, &recv)
		if err != nil {
			logrus.Error("json unmarshal:", err)
			return
		}

		b, err := util.GetEvent(recv.Count)
		if err != nil {
			logrus.Errorf("failed to generate binary data: %v", err)
			return
		}

		c.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			logrus.Error("write:", err)
			return
		}
	}
}
