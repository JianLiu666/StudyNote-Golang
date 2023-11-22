package kvstore

import "context"

type KvStore interface {
	Shutdown(ctx context.Context)
}
