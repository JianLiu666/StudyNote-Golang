package list

import (
	"context"

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

	nextPageKey, err := l.kvstore.GetListHead(context.TODO(), listKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	resp := getListResponse{
		NextPageKey: nextPageKey,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
