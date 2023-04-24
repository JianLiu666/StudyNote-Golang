package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})
}

func main() {
	h2s := &http2.Server{}
	h1s := http.Server{
		Addr:    "localhost:8080",
		Handler: h2c.NewHandler(http.HandlerFunc(handler), h2s),
	}

	err := h1s.ListenAndServeTLS("./certificate.pem", "./private.pem")
	if err != nil {
		logrus.Errorln(err)
	}
}

func handler(resp http.ResponseWriter, req *http.Request) {
	logrus.Infof("Request Protocol: %v", req.Proto)
}
