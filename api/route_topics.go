package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sungoq/sungoq/constants"
)

type CreateTopicReq struct {
	Name string `json:"name"`
}

func (api *API) CreateTopics(c *fiber.Ctx) error {
	input := CreateTopicReq{}
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	err := api.service.Topic.Create(input.Name)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "created",
		"topic":   input.Name,
	})
}

func (api *API) GetAllTopics(c *fiber.Ctx) error {

	topics, err := api.service.Topic.GetAll()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"topics": topics,
	})
}

func (api *API) DeleteTopics(c *fiber.Ctx) error {
	name := c.Query("name", "")
	if name == "" {
		return constants.ErrNameIsEmpty
	}

	err := api.service.Topic.Delete(name)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "deleted",
		"topic":   name,
	})
}
