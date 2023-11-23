package kvstore

import (
	"context"
	"interview20231116/model"
)

type KvStore interface {
	Shutdown(ctx context.Context)

	SetPageToListHead(ctx context.Context, listKey string, page *model.Page) error

	GetHead(listKey string) string

	GetPage(pageKey string)
}
