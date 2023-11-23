package kvstore

import (
	"context"
	"interview20231116/model"
)

type KvStore interface {
	Shutdown(ctx context.Context)

	SetPageToListHead(ctx context.Context, listKey string, page *model.Page) error

	GetListHead(ctx context.Context, listKey string) (string, error)

	GetPage(ctx context.Context, pageKey string) (*model.Page, error)
}
