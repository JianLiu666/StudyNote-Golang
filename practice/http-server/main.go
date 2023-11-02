package main

import (
	"fmt"
	"httpserver/pkg/setting"
	"httpserver/routers"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        routers.InitRouter(),
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
