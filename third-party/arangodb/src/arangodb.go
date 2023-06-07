package src

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	netHttp "net/http"
	"sync"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	HttpRequestTimeout            int = 30
	HttpMaxIdleConnections        int = 60000
	HttpMaxIdleConnectionsPerHost int = 10000
	HttpIdleConnTimeout           int = 150
	ArangoConnectionLimit         int = 100
)

var arangoClient driver.Client
var once sync.Once
var initialized bool

func newArangoClient() {
	once.Do(func() {
		type config struct {
			Addrs    []string `json:"addrs"`
			User     string   `json:"user"`
			Passowrd string   `json:"password"`
		}

		// read configuration file
		content, err := ioutil.ReadFile("./config.json")
		if err != nil {
			panic(err)
		}

		// unmarshal to golang structure
		var conf config
		err = json.Unmarshal(content, &conf)
		if err != nil {
			panic(err)
		}

		dialer := &net.Dialer{
			Timeout:   time.Duration(HttpRequestTimeout) * time.Second,
			KeepAlive: time.Duration(HttpRequestTimeout) * time.Second,
		}

		transport := &netHttp.Transport{
			Proxy:                 netHttp.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConns:          HttpMaxIdleConnections,
			MaxIdleConnsPerHost:   HttpMaxIdleConnectionsPerHost,
			IdleConnTimeout:       time.Duration(HttpIdleConnTimeout) * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		// 注意 ArangoDB 預設對單一個節點的連線數量為32條, 設定成 -1 表示無限制
		// 但必須自行注意機器在單一 Process 上的 fd 大小限制
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: conf.Addrs,
			ConnLimit: ArangoConnectionLimit,
			Transport: transport,
		})
		if err != nil {
			panic(err)
		}

		c, err := driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(conf.User, conf.Passowrd),
		})
		if err != nil {
			panic(err)
		}

		arangoClient = c
		initialized = true
	})
}

func GetDBInstance(ctx context.Context, dbName string) driver.Database {
	if !initialized {
		newArangoClient()
	}

	db, err := arangoClient.Database(ctx, dbName)
	if err != nil {
		panic(err)
	}

	return db
}
