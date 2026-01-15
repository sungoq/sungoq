package api

import (
	"github.com/gofiber/websocket/v2"
	"github.com/sungoq/sungoq/service"
)

func (api *API) Consume(c *websocket.Conn) {
	topic := c.Query("topic", "")
	if topic == "" {
		_ = c.Close()
		return
	}

	messages, err := service.Topic.GetAllMessages(topic)
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

			_ = service.Topic.DeleteMessage(topic, m.ID)
		}
	}()

	for pub := range api.chPublishing {
		if topic == pub.Topic {
			_ = c.WriteMessage(websocket.TextMessage, pub.Message.ToJSON())
			_ = service.Topic.DeleteMessage(topic, pub.Message.ID)
		}
	}

}
