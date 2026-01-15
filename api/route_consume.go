package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/gowok/gowok/web/request"
	"github.com/gowok/gowok/web/response"
	"github.com/sungoq/sungoq/constants"
	"github.com/sungoq/sungoq/service"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetConsume(w http.ResponseWriter, r *http.Request) {
	topic := request.New(r).Query("topic", "")
	if topic == "" {
		response.New(w).BadRequest(constants.ErrNameIsEmpty)
		return
	}

	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		response.New(w).BadRequest(err)
		return
	}

	messages, err := service.TopicGetAllMessages(topic)
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

			_ = service.TopicDeleteMessage(topic, m.ID)
		}
	}()

	for pub := range chPublishing {
		if topic == pub.Topic {
			_ = c.WriteMessage(websocket.TextMessage, pub.Message.ToJSON())
			_ = service.TopicDeleteMessage(topic, pub.Message.ID)
		}
	}

	fmt.Println(123)

}
