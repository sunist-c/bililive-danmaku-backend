package handler

import "github.com/sunist-c/bililive-danmaku-backend/common/threading"

type messageHandler struct {
	scheduler *MessageHandler
	handler   func(pool *Pool, exit chan struct{})
}

func (m messageHandler) Execute(exit chan struct{}) {
	m.handler(m.scheduler.pool, exit)
}

func (m messageHandler) Stop() (success bool) {
	return true
}

type MessageHandler struct {
	pool    *Pool
	handler *threading.Goroutine
}

func (m *MessageHandler) GetMessageChan() *Pool {
	return m.pool
}

func (m *MessageHandler) Close() (success bool) {
	return m.handler.Close()
}

func (m *MessageHandler) Serve() (success bool) {
	return m.handler.Serve()
}

func NewMessageHandler(function func(pool *Pool, exit chan struct{})) *MessageHandler {
	scheduler := &MessageHandler{
		pool: nil,
	}
	scheduler.handler = threading.NewGoroutine(&messageHandler{
		scheduler: scheduler,
		handler:   function,
	})

	return scheduler
}

func NewMessageHandlerWithMessageChan(messageChan *Pool, function func(pool *Pool, exit chan struct{})) *MessageHandler {
	scheduler := NewMessageHandler(function)
	scheduler.pool = messageChan

	return scheduler
}
