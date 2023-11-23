package page

import (
	"context"
	"interview20231116/model"

	"github.com/gofiber/fiber/v2"
)

type getPageResponse struct {
	Articles    []*model.Article `json:"articles"`
	NextPageKey string           `json:"nextPageKey"`
}

func (p *pageRouter) getPage(c *fiber.Ctx) error {
	pageKey := c.Query("pageKey")

	if pageKey == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	page, err := p.kvstore.GetPage(context.TODO(), pageKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	resp := getPageResponse{
		Articles:    page.Articles,
		NextPageKey: page.NextPageKey,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
