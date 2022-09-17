package websocket

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sunist-c/bililive-danmaku/model"
	"github.com/sunist-c/bililive-danmaku/model/info"
	"github.com/sunist-c/bililive-danmaku/model/message"
	"github.com/sunist-c/bililive-danmaku/model/pool"
)

type Client struct {
	roomInfo      *info.Room
	requestInfo   *info.Request
	connection    *websocket.Conn
	messageChan   *pool.Pool
	heartbeatExit chan struct{}
	Connected     bool
}

func (c *Client) GetRoomInfo() info.Room {
	return *c.roomInfo
}

func (c *Client) SendMessage(packageLength uint32, magic uint16, version uint16, typeID uint32, param uint32, data []byte) (err error) {
	if packageLength == 0 {
		packageLength = uint32(len(data) + 16)
	}

	// 将包的头部信息以大端序方式写入字节数组
	packetHead := new(bytes.Buffer)
	var packageData = []interface{}{
		packageLength,
		magic,
		version,
		typeID,
		param,
	}
	for _, v := range packageData {
		if err = binary.Write(packetHead, binary.BigEndian, v); err != nil {
			log.Printf("writing data to websocket package error: %v\n", err)
			return
		}
	}

	// 将包内数据部分追加到数据包内
	sendData := append(packetHead.Bytes(), data...)

	if err = c.connection.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		log.Printf("send websocket data error: %v\n", err)
		return
	}

	return
}

func (c *Client) ReceivedMessage() {
	index := 0
	for {
		if index >= 5 {
			log.Printf("unknown error, try to reconnect...\n")
			break
		}

		_, msg, err := c.connection.ReadMessage()
		if err != nil {
			log.Printf("reading websocket message error: %v\n", err)
			index++
			continue
		}

		switch msg[11] {
		case 8:
			log.Printf("handshake received, connected successfully\n")
			c.Connected = true
		case 3:
			onlineNow := ByteArrToDecimal(msg[16:])
			if uint32(onlineNow) != c.roomInfo.Online {
				c.roomInfo.Online = uint32(onlineNow)
				log.Printf("popularity index changed: %v\n", uint32(onlineNow))
			}
		case 5:
			if inflated, e := ZlibInflate(msg[16:]); e != nil {
				c.messageChan.Unknown <- msg[16:]
			} else {
				for len(inflated) > 0 {
					length := ByteArrToDecimal(inflated[:4])
					command := model.Json.Get(inflated[16:length], "cmd").ToString()
					switch message.Command(command) {
					case message.CommandNewDanmaku:
						c.messageChan.Danmaku <- inflated[16:length]
					case message.CommandNewGift:
						c.messageChan.Gift <- inflated[16:length]
					case message.CommandWelcome:
						c.messageChan.Audience <- inflated[16:length]
					case message.CommandWelcomeGuard:
						c.messageChan.Guard <- inflated[16:length]
					case message.CommandWelcomeVip:
						c.messageChan.Master <- inflated[16:length]
					}
					inflated = inflated[length:]
				}
			}
		}
	}

	go c.restart(c.roomInfo.RoomID)
}

func (c *Client) Serve() bool {
	u := url.URL{Scheme: "wss", Host: model.ApiDanMuServer.ToString(), Path: "/sub"}
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("connect to %v error: %v\n", u.String(), err)
		return false
	} else {
		c.connection = connection
	}

	log.Printf("connected to %v, living: %v\n", c.roomInfo.RoomID, c.roomInfo.LiveStatus)

	payload, err := json.Marshal(c.requestInfo)
	if err != nil {
		log.Printf("marshal handshake payload error: %v\n", err)
		return false
	}
	log.Printf("connected to %v succes, sending handshake: %v\n", model.ApiDanMuServer, string(payload))

	if err = c.SendMessage(0, 16, 1, 7, 1, payload); err != nil {
		log.Printf("send handshake package error: %v\n", err)
		return false
	}
	go c.Heartbeat()
	go c.ReceivedMessage()

	return true
}

func (c *Client) Heartbeat() {
	for {
		select {
		case <-c.heartbeatExit:
			goto exit
		default:
			if c.Connected {
				obj := []byte("5b6f626a656374204f626a6563745d")
				if err := c.SendMessage(31, 16, 1, 2, 1, obj); err != nil {
					log.Println("heart beat err: ", err)
					continue
				}
				time.Sleep(30 * time.Second)
			}
		}
	}

exit:
	log.Printf("heartbeat exited\n")
}

func (c *Client) Stop() {
	c.heartbeatExit <- struct{}{}
	if c.connection != nil {
		_ = c.connection.Close()
	}
}

func (c *Client) init(realRoomID uint32) {
	c.roomInfo = getRoomInfo(realRoomID)
	c.requestInfo = getRequestInfo(realRoomID)
	c.heartbeatExit = make(chan struct{})
	if c.connection != nil {
		_ = c.connection.Close()
	}
	c.connection = nil
	c.Connected = false
}

func (c *Client) restart(realRoomID uint32) {
	c.heartbeatExit <- struct{}{}
	c.init(realRoomID)
	success := false
	for i := 0; i < 10; i++ {
		success = c.Serve()
		if success {
			log.Printf("successfully restart websocket connection to room %v\n", c.roomInfo.RoomID)
			return
		}
	}

	log.Printf("failed restart websocket connection to room %v, exit\n", c.roomInfo.RoomID)
}

func NewClientWithHandler(realRoomID uint32, handler func(pool *pool.Pool)) *Client {
	client := &Client{}
	client.init(realRoomID)
	client.messageChan = pool.NewPoolWithHandler(handler)
	return client
}

func NewClientWithMessageChan(roomID uint32, p *pool.Pool) *Client {
	if roomID < 1000 {
		roomID, _ = getRealRoomID(uint(roomID))
	}

	client := &Client{}
	client.init(roomID)
	client.messageChan = p
	return client
}
