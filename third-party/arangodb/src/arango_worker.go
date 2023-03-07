package src

import (
	"context"
	"errors"
	"net"
	netHttp "net/http"
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

type ArangoWorkerImp struct {
	client driver.Client
	Db     map[string]driver.Database
}

func NewArangoWorkerImp(addrList []string) *ArangoWorkerImp {
	arangoWorker := new(ArangoWorkerImp)

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
		Endpoints: addrList,
		ConnLimit: ArangoConnectionLimit,
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
	arangoWorker.client = c
	arangoWorker.Db = map[string]driver.Database{}

	return arangoWorker
}

func (imp *ArangoWorkerImp) Version(ctx context.Context) string {
	result, err := imp.client.Version(ctx)
	if err != nil {
		panic(err)
	}

	return result.String()
}

func (imp *ArangoWorkerImp) CacheDatabase(ctx context.Context, dbName string) bool {
	db, err := imp.client.Database(ctx, dbName)
	if err != nil {
		panic(err)
	}

	imp.Db[dbName] = db
	return true
}

func (imp *ArangoWorkerImp) GetDatabase(dbName string) (driver.Database, error) {
	db, exists := imp.Db[dbName]
	if exists {
		return db, nil
	}

	return nil, errors.New("database not found")
}
