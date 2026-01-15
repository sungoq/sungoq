package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string `json:"id"`
	Body      any    `json:"body"`
	Timestamp uint   `json:"timestamp"`
	Attempts  uint   `json:"attempts"`
}

type Messages []Message

func NewMessage(body any) Message {
	return Message{
		ID:        uuid.NewString(),
		Body:      body,
		Timestamp: uint(time.Now().Unix()),
	}
}

func (m *Message) ToJSON() []byte {
	messageJson, _ := json.Marshal(m)
	return messageJson
}

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i, j int) bool {
	return m[i].Timestamp < m[j].Timestamp
}

func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
