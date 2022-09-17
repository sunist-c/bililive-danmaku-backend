package message

import "github.com/sunist-c/bililive-danmaku/model"

type Gift struct {
	UserName string `json:"u_uname"`
	Action   string `json:"action"`
	GiftName string `json:"gift_name"`
	Number   uint32 `json:"number"`
	Price    uint32 `json:"price"`
}

func NewGift() *Gift {
	return &Gift{
		UserName: "",
		Action:   "",
		Price:    0,
		GiftName: "",
	}
}

func NewGiftWithData(source []byte) *Gift {
	nums := model.Json.Get(source, "data", "num").ToUint32()
	return &Gift{
		Number:   nums,
		UserName: model.Json.Get(source, "data", "uname").ToString(),
		Action:   model.Json.Get(source, "data", "action").ToString(),
		Price:    model.Json.Get(source, "data", "price").ToUint32() * nums,
		GiftName: model.Json.Get(source, "data", "giftName").ToString(),
	}
}
