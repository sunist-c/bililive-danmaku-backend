package websocket

import (
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
	"github.com/sunist-c/bililive-danmaku-backend/common/threading"
	"github.com/sunist-c/bililive-danmaku-backend/handler"
	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/info"
)

type Client struct {
	roomInfo       *info.Room
	requestInfo    *info.Request
	connection     *websocket.Conn
	messageChan    *handler.Pool
	writer         *writer
	receiver       *receiver
	writerThread   *threading.Goroutine
	receiverThread *threading.Goroutine
	errors         uint
	Connected      bool
}

func (c *Client) initializeWebsocketConnection() error {
	if c.connection != nil {
		_ = c.connection.Close()
	}
	c.connection = nil
	c.Connected = false

	u := url.URL{Scheme: "wss", Host: model.ApiDanMuServer.ToString(), Path: "/sub"}
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logging.Warn("connect to %v error: %v", u.String(), err)
		return err
	} else {
		c.connection = connection
		return nil
	}
}

func (c *Client) tryToConnect() (success bool) {
	for i := 0; i < 3; i++ {
		logging.Info("try to connect to %v, retry %v times", c.roomInfo.RoomID, i)

		payload, err := json.Marshal(c.requestInfo)
		if err != nil {
			logging.Warn("marshal handshake payload error: %v", err)
			continue
		}

		logging.Info("sending handshake package to %v...", model.ApiDanMuServer)

		if err = c.writer.sendMessage(0, 16, 1, 7, 1, payload); err != nil {
			logging.Warn("send handshake package error: %v", err)
			continue
		}

		return true
	}

	return false
}

func (c *Client) restart() {
	c.Close()
	c.Serve()
}

func (c *Client) handleError() {
	c.errors += 1
	if c.errors >= 3 {
		logging.Error("unknown error occurred, try to restart websocket client")
		c.restart()
		c.errors = 0
	}
}

func (c *Client) GetRoomInfo() info.Room {
	return *c.roomInfo
}

func (c *Client) GetMessageChan() *handler.Pool {
	return c.messageChan
}

func (c *Client) Close() (success bool) {
	return c.writerThread.Close() && c.receiverThread.Close()
}

func (c *Client) Serve() (success bool) {
	err := c.initializeWebsocketConnection()
	if err != nil {
		logging.Warn("failed to initialize websocket connection: %v, exit", err)
		return
	}

	if !c.tryToConnect() {
		logging.Warn("failed to connect to websocket server: %v, exit", err)
		return
	}

	if c.writerThread.Serve() && c.receiverThread.Serve() {
		logging.Info("successfully connected to live room: %v", c.roomInfo.RoomID)
	}
	return c.writerThread.Serve() && c.receiverThread.Serve()
}

func NewClient(roomID uint32) *Client {
	if roomID < 1000 {
		realRoomID, err := getRealRoomID(uint(roomID))
		if err != nil {
			logging.Error("get real room ID failed: %v", err)
			return nil
		} else {
			roomID = realRoomID
		}
	}

	client := &Client{
		roomInfo:    getRoomInfo(roomID),
		requestInfo: getRequestInfo(roomID),
		connection:  nil,
		messageChan: nil,
		errors:      0,
		Connected:   false,
	}
	client.writer = newWriter(client)
	client.receiver = newReceiver(client)
	client.writerThread = threading.NewGoroutine(client.writer)
	client.receiverThread = threading.NewGoroutine(client.receiver)

	return client
}

func NewClientWithMessageChan(roomID uint32, pool *handler.Pool) *Client {
	client := NewClient(roomID)
	client.messageChan = pool

	return client
}
