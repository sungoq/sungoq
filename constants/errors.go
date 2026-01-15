package constants

import "errors"

var (
	ErrServiceIsEmpty = errors.New("service is empty")
	ErrNameIsEmpty    = errors.New("name is empty")
	ErrQueueFull      = errors.New("queue is full (maximum size exceeded)")
)
