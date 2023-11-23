package list

import (
	"interview20231116/api/router"
	"interview20231116/pkg/kvstore"

	"github.com/gofiber/fiber/v2"
)

type listRouter struct {
	kvstore kvstore.KvStore
}

func NewListRouter(kvstore kvstore.KvStore) router.Router {
	return &listRouter{
		kvstore: kvstore,
	}
}

func (l *listRouter) Init(r fiber.Router) {
	v1 := r.Group("/v1")

	v1.Post("/list", l.setList)
	v1.Get("/list", l.getList)
}
