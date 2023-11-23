package list

import (
	"context"
	"interview20231116/model"
	"interview20231116/pkg/e"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type setListPayload struct {
	ListKey  string           `json:"listKey"`
	Articles []*model.Article `json:"articles"`
}

func (l *listRouter) setList(c *fiber.Ctx) error {
	payload := &setListPayload{}

	if err := c.BodyParser(payload); err != nil {
		logrus.Errorf("Failed to execute c.BodyParser: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	page := &model.Page{
		Articles: payload.Articles,
	}
	if code := l.kvstore.SetPageToListHead(context.TODO(), payload.ListKey, page); code != e.SUCCESS {
		return c.Status(fiber.StatusBadRequest).SendString(e.GetMsg(code))
	}

	return c.SendStatus(fiber.StatusOK)
}
