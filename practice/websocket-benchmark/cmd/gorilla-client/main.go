package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
	"websocketbenchmark/config"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var bar *progressbar.ProgressBar
var conf *config.Config

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})

	viper.SetConfigFile("./conf.d/env.yaml")
	viper.AutomaticEnv()
	conf = config.NewFromViper()

	opts := progressbar.OptionUseANSICodes(true)
	bar = progressbar.NewOptions64(int64(conf.Simulation.NumClients*conf.Simulation.NumMessages), opts)
}

func main() {
	var clients []*client = make([]*client, conf.Simulation.NumClients)
	var wg sync.WaitGroup
	wg.Add(conf.Simulation.NumClients)

	logrus.Infof("num of clients: %v, each client will send %v messages", conf.Simulation.NumClients, conf.Simulation.NumMessages)

	logrus.Info("start to create clients")

	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%s", conf.Server.Addr, conf.Server.Port),
		Path:   "/echo",
	}
	for i := 0; i < conf.Simulation.NumClients; i++ {
		clients[i] = newClient(u)
		time.Sleep(1 * time.Millisecond)
	}

	logrus.Info("start to send messages")

	for i := 0; i < conf.Simulation.NumClients; i++ {
		go clients[i].start(&wg)
	}

	wg.Wait()

	fmt.Println()
	logrus.Info("start to close connections")
	for i := 0; i < conf.Simulation.NumClients; i++ {
		clients[i].close()
	}

	calculate(clients)
}
