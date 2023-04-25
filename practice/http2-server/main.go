package main

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/sirupsen/logrus"
)

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})
}

func main() {
	h1s := http.Server{
		Addr: "localhost:8080",
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/img", img)

	err := h1s.ListenAndServeTLS("./certificate.pem", "./private.pem")
	if err != nil {
		logrus.Errorln(err)
	}
}

func home(resp http.ResponseWriter, req *http.Request) {
	logrus.Infof("Request Protocol: %v", req.Proto)
}

func img(resp http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadFile("./doge.png")
	if err != nil {
		logrus.Fatalln(err)
	}

	resp.Header().Set("Content-Type", "image/png")

	if strings.Contains(req.Header.Get("accept-encoding"), "br") {
		bw := brotli.NewWriter(resp)
		defer bw.Close()
		resp.Header().Set("Content-Encoding", "br")
		bw.Write(buf)

	} else if strings.Contains(req.Header.Get("accept-encoding"), "gzip") {
		gw := gzip.NewWriter(resp)
		defer gw.Close()
		resp.Header().Set("Content-Encoding", "gzip")
		gw.Write(buf)

	} else {
		resp.Write(buf)
	}
}
