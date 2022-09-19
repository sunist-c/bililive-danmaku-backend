package backend

import (
	"github.com/sunist-c/bililive-danmaku-backend/callback"
	"github.com/sunist-c/bililive-danmaku-backend/websocket"
	"log"
)

type MessageService struct {
	WebsocketClient *websocket.Client
	CallbackClient  *callback.Client
}

func (s *MessageService) Serve() {
	for i := 0; i < 3; i++ {
		if s.WebsocketClient.Serve() {
			return
		} else {
			log.Printf("failed to start websocket service, retry %v times\n", i)
		}
	}

	room := s.WebsocketClient.GetRoomInfo().RoomID
	GetDefaultChannelService().RemoveChannel(room)
	log.Printf("failed to start websocket service, clear up room %v\n", room)
}

func (s *MessageService) Stop() {
	ctx := callback.NewContext("", s.CallbackClient.Options.ExitRouter)
	s.CallbackClient.SendMessage(ctx)
	s.CallbackClient = nil
	s.WebsocketClient.Stop()
	s.WebsocketClient = nil
}

func NewMessageService() *MessageService {
	return &MessageService{}
}
