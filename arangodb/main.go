package main

import (
	"context"
	"log"
)

var arangoClusters []string = []string{"http://192.168.51.118:8551", "http://192.168.51.118:8552", "http://192.168.51.118:8553"}

func main() {
	ctx := context.Background()

	arangodb := NewArangoWorkerImp(ctx, arangoClusters)
	log.Print(arangodb.Version(ctx))

	isCache := arangodb.CacheDatabase(ctx, "Database")
	log.Print(isCache)

	docKey := arangodb.ExplainSave(ctx, "Database", "Demo")
	log.Print(docKey)

	result := arangodb.ExplainTransation(ctx, "Database", "Demo")
	log.Print(result)
}
