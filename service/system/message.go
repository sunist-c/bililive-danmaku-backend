package system

import (
	"fmt"
	"strings"

	"github.com/sunist-c/bililive-danmaku-backend/callback"
	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
	"github.com/sunist-c/bililive-danmaku-backend/handler"
	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/message"
	"github.com/sunist-c/bililive-danmaku-backend/websocket"
)

type MessageService struct {
	callbackClient  *callback.Client
	websocketClient *websocket.Client
	messageClient   *handler.MessageHandler
	initialized     bool
}

func (m MessageService) Close() (success bool) {
	ctx := callback.NewContext("", m.callbackClient.Options.ExitRouter)
	for i := 0; i < 3; i++ {
		if m.callbackClient.SendMessage(ctx) {
			break
		}
		logging.Warn("cannot send exit service message to %v", m.callbackClient.Options.BaseRouter)
	}
	return m.websocketClient.Close() && m.messageClient.Close() && m.callbackClient.Close()
}

func (m MessageService) Serve() (success bool) {
	if !m.initialized {
		logging.Error("cannot serve an uninitialized service")
		return false
	}

	return m.websocketClient.Serve() && m.messageClient.Serve() && m.callbackClient.Serve()
}

func (m *MessageService) InitializeMessageService(roomID uint32, option *callback.ClientOptions, implement func(pool *handler.Pool, exit chan struct{})) *MessageService {
	messageChan := handler.NewPool()
	m.callbackClient = callback.NewClientWithOption(option)
	m.messageClient = handler.NewMessageHandlerWithMessageChan(messageChan, implement)
	m.websocketClient = websocket.NewClientWithMessageChan(roomID, messageChan)
	m.initialized = true

	return m
}

func NewMessageService() *MessageService {
	return &MessageService{
		initialized:     false,
		callbackClient:  nil,
		messageClient:   nil,
		websocketClient: nil,
	}
}

func WebhookImplementation(messageService *MessageService) func(pool *handler.Pool, exit chan struct{}) {
	return func(pool *handler.Pool, exit chan struct{}) {
		for {
			select {
			case <-exit:
				return
			case uc := <-pool.Unknown:
				if cmd := model.Json.Get(uc, "cmd").ToString(); message.Command(cmd) == message.CommandRoomFocusedChange {
					fans := model.Json.Get(uc, "data", "fans").ToInt()
					customMessage := message.NewCustom(fmt.Sprintf("room fans changed: %v", fans))
					messageService.callbackClient.SendMessage(callback.NewContext(&customMessage, messageService.callbackClient.Options.FansRouter))
				}
			case src := <-pool.Danmaku:
				m := message.NewDanmakuWithData(src)
				messageService.callbackClient.SendMessage(callback.NewContext(m, messageService.callbackClient.Options.DanmakuRouter))
			case src := <-pool.Gift:
				g := message.NewGiftWithData(src)
				messageService.callbackClient.SendMessage(callback.NewContext(g, messageService.callbackClient.Options.GiftRouter))
			case src := <-pool.Audience:
				name := model.Json.Get(src, "data", "uname").ToString()
				welcomeMessage := message.NewWelcome(name, "audience")
				messageService.callbackClient.SendMessage(callback.NewContext(&welcomeMessage, messageService.callbackClient.Options.AudienceRouter))
			case src := <-pool.Guard:
				name := model.Json.Get(src, "data", "username").ToString()
				welcomeMessage := message.NewWelcome(name, "guard")
				messageService.callbackClient.SendMessage(callback.NewContext(&welcomeMessage, messageService.callbackClient.Options.GuardRouter))
			case src := <-pool.Master:
				cw := model.Json.Get(src, "data", "copy_writing").ToString()
				cw = strings.Replace(cw, "<%", "", 1)
				cw = strings.Replace(cw, "%>", "", 1)
				customMessage := message.NewCustom(cw)
				messageService.callbackClient.SendMessage(callback.NewContext(&customMessage, messageService.callbackClient.Options.MasterRouter))
			}
		}
	}
}
