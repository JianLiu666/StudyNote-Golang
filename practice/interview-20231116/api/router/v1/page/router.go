package page

import (
	"interview20231116/api/router"
	"interview20231116/pkg/kvstore"

	"github.com/gofiber/fiber/v2"
)

type pageRouter struct {
	kvstore kvstore.KvStore
}

func NewPageRouter(kvstore kvstore.KvStore) router.Router {
	return &pageRouter{
		kvstore: kvstore,
	}
}

func (p *pageRouter) Init(r fiber.Router) {
	v1 := r.Group("/v1")
	v1.Get("/page", p.getPage)

}
