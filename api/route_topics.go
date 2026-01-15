package api

import (
	"net/http"

	"github.com/gowok/gowok/web/request"
	"github.com/gowok/gowok/web/response"
	"github.com/ngamux/ngamux"
	"github.com/sungoq/sungoq/constants"
	"github.com/sungoq/sungoq/service"
)

type CreateTopicReq struct {
	Name string `json:"name"`
}

func PostTopics(w http.ResponseWriter, r *http.Request) {
	res := response.New(w)

	input := CreateTopicReq{}
	if err := request.New(r).JSON(&input); err != nil {
		res.BadRequest(err)
		return
	}

	err := service.TopicCreate(input.Name)
	if err != nil {
		res.BadRequest(err)
		return
	}

	response.New(w).JSON(ngamux.Map{
		"message": "created",
		"topic":   input.Name,
	})
}

func GetTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := service.TopicGetAll()
	if err != nil {
		response.New(w).BadRequest(err)
		return
	}

	response.New(w).JSON(ngamux.Map{
		"topics": topics,
	})
}

func DeleteTopics(w http.ResponseWriter, r *http.Request) {
	name := request.New(r).Query("name")
	if name == "" {
		response.New(w).BadRequest(constants.ErrNameIsEmpty)
		return
	}

	err := service.TopicDelete(name)
	if err != nil {
		response.New(w).BadRequest(err)
		return
	}

	response.New(w).JSON(ngamux.Map{
		"message": "deleted",
		"topic":   name,
	})
}
