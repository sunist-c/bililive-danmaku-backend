package model

import (
	"sync"

	iter "github.com/json-iterator/go"
)

var (
	Json      = iter.ConfigCompatibleWithStandardLibrary
	config    *Config
	WaitGroup = &sync.WaitGroup{}
)

const (
	ApiGetRealRoomID  Api = "http://api.live.bilibili.com/room/v1/Room/room_init"                 // params: id=xxx
	ApiDanMuServer    Api = "broadcastlv.chat.bilibili.com:443"                                   // params: null
	ApiGetAccessToken Api = "https://api.live.bilibili.com/room/v1/Danmu/getConf"                 // params: room_id=xxx&platform=pc&player=web
	ApiGetRoomInfo    Api = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom" // params: room_id=xxx
)

func GetGlobalConfig() *Config {
	return config
}
