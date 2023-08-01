package main

import (
	"arangodb/src"
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
)

var dbName string = "Test"

func main() {
	ctx := context.Background()

	db := src.GetDBInstance(ctx, dbName)

	for _, j := range jobs {
		fmt.Println(j.Name)
		cursor, err := db.Query(driver.WithQueryProfile(ctx), j.AQL, j.Binds)
		if err != nil {
			panic(err)
		}
		defer cursor.Close()
		fmt.Println(float64(cursor.Extra().GetStatistics().ExecutionTime().Milliseconds()) / 1000)
	}
}
