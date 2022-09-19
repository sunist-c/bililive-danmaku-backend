package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/bililive-danmaku-backend/websocket"
	"log"
	"strings"

	"github.com/sunist-c/bililive-danmaku-backend/callback"
	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/message"
	"github.com/sunist-c/bililive-danmaku-backend/model/pool"
)

type RegisterRequest struct {
	RoomID           uint32 `json:"room_id"`
	BaseRouter       string `json:"baseRouter"`
	DanmakuRouter    string `json:"danmakuRouter"`
	GiftRouter       string `json:"giftRouter"`
	GuardRouter      string `json:"guardRouter"`
	MasterRouter     string `json:"masterRouter"`
	AudienceRouter   string `json:"audienceRouter"`
	FansRouter       string `json:"fansRouter"`
	CustomMessage    string `json:"customMessage"`
	ExitRouter       string `json:"exitRouter"`
	UpdateRoomRouter string `json:"updateRoomRouter"`
}

type RoomInfo struct {
	RoomID     uint32 `json:"room_id"`
	UpUID      uint32 `json:"up_uid"`
	Title      string `json:"title"`
	Online     uint32 `json:"online"`
	Tags       string `json:"tags"`
	LiveStatus bool   `json:"live_status"`
}

type RegisterResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WebsocketHandler(callbackClient *callback.Client) func(pool *pool.Pool) {
	return func(pool *pool.Pool) {
		for {
			select {
			case uc := <-pool.Unknown:
				if cmd := model.Json.Get(uc, "cmd").ToString(); message.Command(cmd) == message.CommandRoomFocusedChange {
					fans := model.Json.Get(uc, "data", "fans").ToInt()
					log.Printf("room fans changed: %v\n", fans)
					customMessage := message.NewCustom(fmt.Sprintf("room fans changed: %v", fans))
					callbackClient.SendMessage(callback.NewContext(&customMessage, callbackClient.Options.FansRouter))
				}
			case src := <-pool.Danmaku:
				m := message.NewDanmakuWithData(src)
				log.Printf("Lv%d%s - Lv%d%s: %s\n", m.MedalLevel, m.MedalName, m.UserLevel, m.UserName, m.Message)
				callbackClient.SendMessage(callback.NewContext(m, callbackClient.Options.DanmakuRouter))
			case src := <-pool.Gift:
				g := message.NewGiftWithData(src)
				log.Printf("%s %s %s*%v, total valued %v\n", g.UserName, g.Action, g.GiftName, g.Number, g.Price)
				callbackClient.SendMessage(callback.NewContext(g, callbackClient.Options.GiftRouter))
			case src := <-pool.Audience:
				name := model.Json.Get(src, "data", "uname").ToString()
				log.Printf("welcome master %s entered room\n", name)
				welcomeMessage := message.NewWelcome(name, "audience")
				callbackClient.SendMessage(callback.NewContext(&welcomeMessage, callbackClient.Options.AudienceRouter))
			case src := <-pool.Guard:
				name := model.Json.Get(src, "data", "username").ToString()
				log.Printf("welcome guard %s entered room\n", name)
				welcomeMessage := message.NewWelcome(name, "guard")
				callbackClient.SendMessage(callback.NewContext(&welcomeMessage, callbackClient.Options.GuardRouter))
			case src := <-pool.Master:
				cw := model.Json.Get(src, "data", "copy_writing").ToString()
				cw = strings.Replace(cw, "<%", "", 1)
				cw = strings.Replace(cw, "%>", "", 1)
				log.Printf("%s\n", cw)
				customMessage := message.NewCustom(cw)
				callbackClient.SendMessage(callback.NewContext(&customMessage, callbackClient.Options.MasterRouter))
			}
		}
	}
}

func RegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &RegisterRequest{}
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			ctx.JSON(400, RegisterResponse{Success: false, Error: err.Error()})
			return
		}

		option := &callback.ClientOptions{
			BaseRouter:       request.BaseRouter,
			DanmakuRouter:    request.DanmakuRouter,
			GiftRouter:       request.GiftRouter,
			GuardRouter:      request.GuardRouter,
			MasterRouter:     request.MasterRouter,
			AudienceRouter:   request.AudienceRouter,
			FansRouter:       request.FansRouter,
			CustomMessage:    request.CustomMessage,
			ExitRouter:       request.ExitRouter,
			UpdateRoomRouter: request.UpdateRoomRouter,
		}

		service := NewMessageService()
		service.CallbackClient = callback.NewClientWithOption(option)
		messageChan := pool.NewPoolWithHandler(WebsocketHandler(service.CallbackClient))
		service.WebsocketClient = websocket.NewClientWithMessageChan(request.RoomID, messageChan)

		GetDefaultChannelService().AddChannel(request.RoomID, service)

		service.Serve()

		ctx.JSON(200, "success")
	}
}
