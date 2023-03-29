package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
	"websocketbenchmark/internal/config"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const numClients int = 1     //
const numMessages int = 1000 //

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
	var clients [numClients]*client
	var wg sync.WaitGroup

	conf := config.NewFromViper()

	addr := fmt.Sprintf("%s:%s", conf.Server.Addr, conf.Server.Port)
	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   "/echo",
	}

	wg.Add(numClients)
	for i := 0; i < numClients; i++ {
		clients[i] = newClient(u)
		time.Sleep(1 * time.Millisecond)
	}

	for i := 0; i < numClients; i++ {
		go clients[i].start(&wg)
	}

	wg.Wait()

	for i := 0; i < numClients; i++ {
		clients[i].close()
	}

	calculate(clients)
}
