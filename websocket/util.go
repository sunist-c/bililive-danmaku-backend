package websocket

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/info"
)

func byteArrToInt(src []byte) (sum int) {
	if src == nil {
		return 0
	}
	b := []byte(hex.EncodeToString(src))
	l := len(b)
	for i := l - 1; i >= 0; i-- {
		base := int(math.Pow(16, float64(l-i-1)))
		var mul int
		if int(b[i]) >= 97 {
			mul = int(b[i]) - 87
		} else {
			mul = int(b[i]) - 48
		}

		sum += base * mul
	}

	return sum
}

func zlibInflate(compress []byte) ([]byte, error) {
	var out bytes.Buffer
	c := bytes.NewReader(compress)
	r, err := zlib.NewReader(c)
	if err != zlib.ErrChecksum && err != zlib.ErrDictionary && err != zlib.ErrHeader && r != nil {
		_, _ = io.Copy(&out, r)
		if err = r.Close(); err != nil {
			return nil, err
		}
		return out.Bytes(), nil
	}

	return nil, err
}

func getToken(realRoomID uint32) string {
	url := fmt.Sprintf("%s?room_id=%d&platform=pc&player=web", model.ApiGetAccessToken, realRoomID)
	log.Printf("get room token for %v\n", realRoomID)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("get room token error: %v\n", err)
		return ""
	}

	raw, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		log.Printf("read room token response error: %v\n", err)
		return ""
	}
	token := model.Json.Get(raw, "data").Get("token").ToString()

	return token
}

func getRoomInfo(realRoomID uint32) *info.Room {
	url := fmt.Sprintf("%s?room_id=%d", model.ApiGetRoomInfo, realRoomID)
	log.Printf("get room info for %v\n", realRoomID)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("get room info error: %v\n", err)
		return nil
	}

	raw, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		log.Printf("read room info response error: %v\n", err)
		return nil
	}

	roomInfo := &info.Room{}
	roomInfo.RoomID = realRoomID
	roomInfo.UpUID = model.Json.Get(raw, "data").Get("room_info").Get("uid").ToUint32()
	roomInfo.Title = model.Json.Get(raw, "data").Get("room_info").Get("title").ToString()
	roomInfo.Tags = model.Json.Get(raw, "data").Get("room_info").Get("tags").ToString()
	roomInfo.LiveStatus = model.Json.Get(raw, "data").Get("room_info").Get("live_status").ToBool()
	roomInfo.LockStatus = model.Json.Get(raw, "data").Get("room_info").Get("lock_status").ToBool()

	return roomInfo
}

func getRequestInfo(realRoomID uint32) *info.Request {
	token := getToken(realRoomID)
	return &info.Request{
		UID:           0,
		RoomID:        realRoomID,
		ProtoVersion:  2,
		Platform:      "web",
		ClientVersion: "1.10.2",
		Type:          2,
		Key:           token,
	}
}

func getRealRoomID(shortID uint) (realID uint32, err error) {
	url := fmt.Sprintf("%s?id=%d", model.ApiGetRealRoomID, shortID)
	log.Printf("getting the real room id for %v\n", shortID)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("get real room id error: %v\n", err)
		return 0, err
	}

	raw, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		log.Printf("read real room id response error: %v\n", err)
		return 0, err
	}
	realID = model.Json.Get(raw, "data", "room_id").ToUint32()

	return realID, nil
}
