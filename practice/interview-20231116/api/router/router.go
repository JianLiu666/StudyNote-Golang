package router

import "github.com/gofiber/fiber/v2"

type Router interface {
	Init(r fiber.Router)
}
