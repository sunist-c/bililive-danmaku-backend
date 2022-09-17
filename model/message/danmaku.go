package message

import "github.com/sunist-c/bililive-danmaku/model"

type Danmaku struct {
	UID        uint32 `json:"uid"`
	UserName   string `json:"uname"`
	UserLevel  uint32 `json:"ulevel"`
	Message    string `json:"text"`
	MedalLevel uint32 `json:"medal_level"`
	MedalName  string `json:"medal_name"`
}

func NewDanmaku() *Danmaku {
	return &Danmaku{
		UID:        0,
		UserName:   "",
		UserLevel:  0,
		Message:    "",
		MedalLevel: 0,
		MedalName:  "无勋章",
	}
}

func NewDanmakuWithData(source []byte) *Danmaku {
	d := &Danmaku{
		UID:        model.Json.Get(source, "info", 2, 0).ToUint32(),
		UserName:   model.Json.Get(source, "info", 2, 1).ToString(),
		UserLevel:  model.Json.Get(source, "info", 4, 0).ToUint32(),
		Message:    model.Json.Get(source, "info", 1).ToString(),
		MedalName:  model.Json.Get(source, "info", 3, 1).ToString(),
		MedalLevel: model.Json.Get(source, "info", 3, 0).ToUint32(),
	}

	if d.MedalName == "" {
		d.MedalName = "无勋章"
	}

	return d
}
