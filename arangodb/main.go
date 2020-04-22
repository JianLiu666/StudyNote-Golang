package main

import (
	"context"
	"fmt"
	"log"
	"net"
	netHttp "net/http"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	MaxIdleConns        int = 60000
	MaxIdleConnsPerHost int = 10000
	IdleConnTimeout     int = 150
	RequestTimeout      int = 30
)

var arangoClusters []string = []string{"http://192.168.51.118:8551", "http://192.168.51.118:8552", "http://192.168.51.118:8553"}

func main() {
	ctx := context.Background()
	transport := &netHttp.Transport{
		Proxy: netHttp.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(RequestTimeout) * time.Second,
			KeepAlive: time.Duration(RequestTimeout) * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          MaxIdleConns,
		MaxIdleConnsPerHost:   MaxIdleConnsPerHost,
		IdleConnTimeout:       time.Duration(IdleConnTimeout) * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: arangoClusters,
		Transport: transport,
	})
	if err != nil {
		panic(err)
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
		// Authentication: driver.BasicAuthentication("user", "password"),
	})
	if err != nil {
		panic(err)
	}
	log.Println("access ok")
	log.Println(c.Version(ctx))

	db, err := c.Database(ctx, "_system")
	if err != nil {
		panic(err)
	}

	found, err := db.CollectionExists(ctx, "AdminAccounts")
	if err != nil {
		panic(err)
	}
	fmt.Println(found)
}
