package main

import (
	"context"
	"log"
)

var arangoClusters []string = []string{"http://localhost:8529"}

func main() {
	ctx := context.Background()

	arangodb := NewArangoWorkerImp(ctx, arangoClusters)
	log.Print(arangodb.Version(ctx))

	isCache := arangodb.CacheDatabase(ctx, "_system")
	log.Print(isCache)

	docKey := arangodb.ExplainSave(ctx, "_system", "Demo")
	log.Print(docKey)

	result := arangodb.ExplainTransation(ctx, "_system", "Demo")
	log.Print(result)
}
