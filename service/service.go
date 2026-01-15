package service

import "github.com/sungoq/sungoq/service/topic"

type Service struct {
	Topic *topic.TopicService
}

func New() (*Service, error) {
	topicService, err := topic.New()
	if err != nil {
		return nil, err
	}

	return &Service{
		Topic: topicService,
	}, nil
}
