package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"websocketbenchmark/internal/config"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const numClients int = 100    //
const numMessages int = 10000 //

var bar *progressbar.ProgressBar

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})

	viper.SetConfigFile("./conf.d/env.yaml")
	viper.AutomaticEnv()

	opts := progressbar.OptionUseANSICodes(true)
	bar = progressbar.NewOptions64(int64(numClients)*int64(numMessages), opts)
}

func main() {
	conf := config.NewFromViper()

	var clients [numClients]*client
	var wg sync.WaitGroup
	wg.Add(numClients)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	addr := fmt.Sprintf("http://%s:%s", conf.Server.Addr, conf.Server.Port)

	logrus.Infof("num of clients: %v, each client will send %v messages", numClients, numMessages)

	logrus.Info("start to create clients")
	for i := 0; i < numClients; i++ {
		clients[i] = newClient(ctx, addr)
		time.Sleep(1 * time.Millisecond)
	}

	logrus.Info("start to send messages")
	for i := 0; i < numClients; i++ {
		go clients[i].start(ctx, &wg)
	}

	wg.Wait()

	for i := 0; i < numClients; i++ {
		clients[i].close()
	}

	calculate(clients)
}
