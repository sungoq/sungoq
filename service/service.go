package service

import "github.com/sungoq/sungoq/service/topic"

var Topic *topic.TopicService

func init() {
	topicService, err := topic.New()
	if err != nil {
		panic(err)
	}
	Topic = topicService
}
