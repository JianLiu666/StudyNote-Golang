package list

import (
	"context"
	"interview20231116/pkg/e"

	"github.com/gofiber/fiber/v2"
)

type getListResponse struct {
	NextPageKey string `json:"nextPageKey"`
}

func (l *listRouter) getList(c *fiber.Ctx) error {
	listKey := c.Query("listKey")

	if listKey == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	nextPageKey, code := l.kvstore.GetListHead(context.TODO(), listKey)
	if code != e.SUCCESS {
		return c.Status(fiber.StatusBadRequest).SendString(e.GetMsg(code))
	}

	resp := getListResponse{
		NextPageKey: nextPageKey,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
