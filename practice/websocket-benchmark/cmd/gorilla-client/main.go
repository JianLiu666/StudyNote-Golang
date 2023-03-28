package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"
	"websocketbenchmark/internal/config"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	addr := fmt.Sprintf("%s:%s", conf.Server.Addr, conf.Server.Port)
	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   "/echo",
	}
	logrus.Infof("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logrus.Error("read:", err)
				return
			}
			logrus.Info("recv:", string(message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				logrus.Error("write:", err)
				return
			}
		case <-interrupt:
			logrus.Error("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logrus.Error("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
