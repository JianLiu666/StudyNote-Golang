package main

import (
	"context"
	"fmt"
	"time"
	"websocketbenchmark/internal/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	addr := fmt.Sprintf("http://%s:%s", conf.Server.Addr, conf.Server.Port)
	c, _, err := websocket.Dial(ctx, addr, &websocket.DialOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		logrus.Fatal("dial: ", err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for i := 0; i < 5; i++ {
		err = wsjson.Write(ctx, c, map[string]int{
			"i": i,
		})
		if err != nil {
			logrus.Fatal(err)
		}

		v := map[string]int{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info(v)

		if v["i"] != i {
			logrus.Fatalf("expected %v but got %v", i, v)
		}
	}

	c.Close(websocket.StatusNormalClosure, "")
}
