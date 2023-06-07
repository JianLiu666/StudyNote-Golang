package main

import (
	"arangodb/src"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
)

func main() {
	ctx := context.Background()
	instance := src.GetDBInstance(ctx, "_system")

	explainSave(ctx, instance)

	explainTransaction(ctx, instance)
}

func explainSave(ctx context.Context, instance driver.Database) {
	col, err := instance.Collection(ctx, "Demo")
	if err != nil {
		panic(err)
	}

	meta, err := col.CreateDocument(ctx, struct {
		Text       string `json:"Text"`
		CreateTime int64  `json:"CreateTime"`
		Number     int64  `json:"Number"`
	}{
		Text:       "Write something.",
		CreateTime: time.Now().Unix(),
		Number:     50000000000,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(meta.Key)
}

func explainTransaction(ctx context.Context, instance driver.Database) {
	type dao struct {
		Text       string `json:"Text"`
		CreateTime int64  `json:"CreateTime"`
		Number     int64  `json:"Number"`
	}

	doc := dao{
		Text:       "Write something.",
		CreateTime: time.Now().Unix(),
		Number:     50000000000,
	}
	jsonData, err := json.Marshal(&doc)
	if err != nil {
		panic(err)
	}

	// transaction script by javascript
	// creation and deletion databases/collections/indexes are NOT ALLOWED inside transactions
	action := `function (Params) {
		const db = require('@arangodb').db;
		const col  = db._collection(Params[0]);
		const meta = col.save(JSON.parse(Params[1]));
		const result = col.firstExample({
			"_key": meta._key,
		})
		return JSON.stringify(result);
	}`

	options := &driver.TransactionOptions{
		// Transaction store all keys and values in RAM, so this setting ensure that does not happen any
		// out-of-memory situations by limiting the size of any individual transaction.
		MaxTransactionSize: 1000,
		// Those collection will be used in write or read mode.
		WriteCollections: []string{"Demo"},
		// Those collections will be used in read-only mode.
		ReadCollections: []string{"Demo"},
		// Some parameters will be used in transaction script.
		Params: []interface{}{"Demo", string(jsonData)},
		// Whether the transaction is forced to be synchronous.
		WaitForSync: false,
	}

	result, err := instance.Transaction(ctx, action, options)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
