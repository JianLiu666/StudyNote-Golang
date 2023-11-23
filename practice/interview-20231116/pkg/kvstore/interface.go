package kvstore

import (
	"context"
	"interview20231116/model"
	"interview20231116/pkg/e"
)

type KvStore interface {
	Shutdown(ctx context.Context)

	SetPageToListHead(ctx context.Context, listKey string, page *model.Page) e.CODE

	GetListHead(ctx context.Context, listKey string) (string, e.CODE)

	GetPage(ctx context.Context, pageKey string) (*model.Page, e.CODE)
}
