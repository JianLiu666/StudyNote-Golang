package head

import (
	"github.com/gofiber/fiber/v2"
)

func NewRouter(r fiber.Router) {
	v1 := r.Group("/v1")
	v1.Post("/head", SetHead)
	v1.Get("/head", GetHead)
}
