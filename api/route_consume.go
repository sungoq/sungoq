package api

import (
	"github.com/gofiber/websocket/v2"
)

func (api *API) Consume(c *websocket.Conn) {
	topic := c.Query("topic", "")
	if topic == "" {
		_ = c.Close()
		return
	}

	messages, err := api.service.Topic.GetAllMessages(topic)
	if err != nil {
		_ = c.Close()
		return
	}

	go func() {
		for _, m := range messages {
			mJson := m.ToJSON()
			if err := c.WriteMessage(websocket.TextMessage, mJson); err != nil {
				continue
			}

			_ = api.service.Topic.DeleteMessage(topic, m.ID)
		}
	}()

	for pub := range api.chPublishing {
		if topic == pub.Topic {
			_ = c.WriteMessage(websocket.TextMessage, pub.Message.ToJSON())
			_ = api.service.Topic.DeleteMessage(topic, pub.Message.ID)
		}
	}

}
