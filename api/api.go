package api

import (
	"github.com/gowok/gowok"
	"github.com/sungoq/sungoq/model"
)

var chPublishing chan model.Publishing

func Configure() {
	chPublishing = make(chan model.Publishing, 1)
	api := gowok.Web
	api.Post("/topics", PostTopics)
	api.Get("/topics", GetTopics)
	api.Delete("/topics", DeleteTopics)

	api.Post("/publish", PostPublish)

	api.Get("/consume", GetConsume)
}
