package api

import (
	"net/http"

	"github.com/gowok/gowok/web/request"
	"github.com/gowok/gowok/web/response"
	"github.com/ngamux/ngamux"
	"github.com/sungoq/sungoq/model"
	"github.com/sungoq/sungoq/service"
)

func PostPublish(w http.ResponseWriter, r *http.Request) {
	res := response.New(w)
	input := model.Publishing{}
	if err := request.New(r).JSON(&input); err != nil {
		res.BadRequest(err)
		return
	}

	message, err := service.TopicPublish(input.Topic, input.Message)
	if err != nil {
		res.BadRequest(err)
		return
	}

	chPublishing <- model.Publishing{
		Topic:   input.Topic,
		Message: message,
	}

	res.JSON(ngamux.Map{
		"message": "sent",
		"topic":   input.Topic,
		"id":      message.ID,
	})
}
