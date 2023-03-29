package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
	"websocketbenchmark/internal/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type data struct {
	Count     int32 `json:"c"`
	Timestamp int64 `json:"ts"`
}

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

	addr := fmt.Sprintf("%s:%s", conf.Server.Addr, conf.Server.Port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Error("listening: ", err)
		return
	}
	logrus.Infof("listening on http://%v", l.Addr())

	s := &http.Server{
		Handler:      echoServer{},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		logrus.Error("shutdown: ", err)
	}
}
