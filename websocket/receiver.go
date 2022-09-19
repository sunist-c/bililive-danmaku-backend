package websocket

import (
	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/message"
)

type receiver struct {
	client *Client
}

func (r receiver) Execute(exit chan struct{}) {
	for {
		select {
		case <-exit:
			return
		default:
			_, msg, err := r.client.connection.ReadMessage()
			if err != nil {
				logging.Warn("reading websocket message error: %v", err)
				r.client.handleError()
				continue
			}

			switch msg[11] {
			// connection changed message
			case 8:
				logging.Info("handshake received, connected successfully\n")
				r.client.Connected = true

			// online audience count changed message
			case 3:
				onlineNow := byteArrToInt(msg[16:])
				if uint32(onlineNow) != r.client.roomInfo.Online {
					r.client.roomInfo.Online = uint32(onlineNow)
					logging.Info("popularity index changed: %v", uint32(onlineNow))
				}

			// common message
			case 5:
				inflated, e := zlibInflate(msg[16:])

				// uncompressed message
				if e != nil {
					r.client.messageChan.Unknown <- msg[16:]
				}

				// compressed message
				for len(inflated) > 0 {
					length := byteArrToInt(inflated[:4])
					command := model.Json.Get(inflated[16:length], "cmd").ToString()
					switch message.Command(command) {
					// danmaku message
					case message.CommandNewDanmaku:
						r.client.messageChan.Danmaku <- inflated[16:length]

					// gift message
					case message.CommandNewGift:
						r.client.messageChan.Gift <- inflated[16:length]

					// new audience message
					case message.CommandWelcome:
						r.client.messageChan.Audience <- inflated[16:length]

					// guard welcome message
					case message.CommandWelcomeGuard:
						r.client.messageChan.Guard <- inflated[16:length]

					// master welcome message
					case message.CommandWelcomeVip:
						r.client.messageChan.Master <- inflated[16:length]
					}
					inflated = inflated[length:]
				}
			}

			if !r.client.Connected {
				logging.Info("websocket connection closed, exit receiver...")
				exit <- struct{}{}
			}
		}
	}

}

func (r receiver) Stop() (success bool) {
	return true
}

func newReceiver(client *Client) *receiver {
	return &receiver{client: client}
}
