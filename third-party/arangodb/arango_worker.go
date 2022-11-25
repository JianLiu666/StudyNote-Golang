package main

import (
	"context"
	"encoding/json"
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
	ConnectionLimit     int = 100
)

type ArangoWorkerImp struct {
	client driver.Client
	Db     map[string]driver.Database
}

func NewArangoWorkerImp(ctx context.Context, addrList []string) *ArangoWorkerImp {
	arangoWorker := new(ArangoWorkerImp)

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

	// 注意 ArangoDB 預設對單一節點的連線數量為32條, 設定成 -1 表示無限制
	// 但必須自行注意機器在單一 Process 上的 fd 大小限制
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: arangoClusters,
		ConnLimit: ConnectionLimit,
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

func (this *ArangoWorkerImp) Version(ctx context.Context) string {
	result, err := this.client.Version(ctx)
	if err != nil {
		panic(err)
	}

	return result.String()
}

func (this *ArangoWorkerImp) CacheDatabase(ctx context.Context, dbName string) bool {
	db, err := this.client.Database(ctx, dbName)
	if err != nil {
		panic(err)
	}

	this.Db[dbName] = db
	return true
}

/** ArangoDB create document example
 *
 * @param ctx golang context
 * @param dbName name of database
 * @param colName name of collection
 * @return string key of new document */
func (this *ArangoWorkerImp) ExplainSave(ctx context.Context, dbName, colName string) string {
	if _, exists := this.Db[dbName]; !exists {
		return ""
	}

	col, err := this.Db[dbName].Collection(ctx, colName)
	if err != nil {
		panic(err)
	}

	meta, err := col.CreateDocument(ctx, struct {
		Text       string `json:"Text"`
		CreateTime int64  `json:"CreateTime"`
	}{
		Text:       "Write something.",
		CreateTime: time.Now().Unix(),
	})
	if err != nil {
		panic(err)
	}

	return meta.Key
}

/** ArangoDB transation example
 *
 * @param ctx golang context
 * @param dbName name of database
 * @param colName name of collection
 * @return interface{} result transation result */
func (this *ArangoWorkerImp) ExplainTransation(ctx context.Context, dbName, colName string) interface{} {
	if _, exists := this.Db[dbName]; !exists {
		return ""
	}

	doc := struct {
		Text       string `json:"Text"`
		CreateTime int64  `json:"CreateTime"`
	}{
		Text:       "Write something.",
		CreateTime: time.Now().Unix(),
	}
	jsonData, err := json.Marshal(&doc)
	if err != nil {
		panic(err)
	}

	// transation script by javascript
	// creation and deletion databases/collections/indexes are NOT ALLOWED inside transations
	action := `function (Params) {
		const db = require('@arangodb').db;
		const col  = db._collection(Params[0]);
		const result = col.save(JSON.parse(Params[1]));
		return result;}`

	options := &driver.TransactionOptions{
		// Transation store all keys and values in RAM, so this setting ensure that does not happen any
		// out-of-memory situations by limiting the size of any individual transtation.
		MaxTransactionSize: 1000,
		// Those collection will be used in write or read mode.
		WriteCollections: []string{colName},
		// Those collections will be used in read-only mode.
		ReadCollections: []string{colName},
		// Some parameters will be used in trasation script.
		Params: []interface{}{colName, string(jsonData)},
		// Wheher the transation is forced to be synchronous.
		WaitForSync: false,
	}

	result, err := this.Db[dbName].Transaction(context.TODO(), action, options)
	if err != nil {
		panic(err)
	}

	return result
}
