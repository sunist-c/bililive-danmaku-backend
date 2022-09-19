package websocket

import (
	"bytes"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
	"sync"
	"time"
)

type writer struct {
	mu     *sync.Mutex
	client *Client
}

func (w *writer) sendMessage(packageLength uint32, magic uint16, version uint16, typeID uint32, param uint32, data []byte) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
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
			logging.Warn("writing data to websocket package error: %v", err)
			return err
		}
	}

	// 将包内数据部分追加到数据包内
	sendData := append(packetHead.Bytes(), data...)

	if err = w.client.connection.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		logging.Warn("send websocket data error: %v", err)
		return err
	}

	return nil
}

func (w writer) Execute(exit chan struct{}) {
	for {
		select {
		case <-exit:
			return
		default:
			payload := []byte("5b6f626a656374204f626a6563745d")
			if err := w.sendMessage(31, 16, 1, 2, 1, payload); err != nil {
				logging.Warn("heart beat err: %v, retrying...", err)
				continue
			} else {
				logging.Debug("successfully send heart beat package")
			}
			time.Sleep(30 * time.Second)

			if !w.client.Connected {
				logging.Info("websocket connection closed, exit writer...")
				exit <- struct{}{}
			}
		}
	}
}

func (w writer) Stop() (success bool) {
	return true
}

func newWriter(client *Client) *writer {
	return &writer{
		mu:     &sync.Mutex{},
		client: client,
	}
}
