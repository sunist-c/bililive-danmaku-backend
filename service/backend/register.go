package backend

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sunist-c/bililive-danmaku/callback"
	"github.com/sunist-c/bililive-danmaku/model"
	"github.com/sunist-c/bililive-danmaku/model/info"
	"github.com/sunist-c/bililive-danmaku/model/message"
	"github.com/sunist-c/bililive-danmaku/model/pool"
	"github.com/sunist-c/bililive-danmaku/websocket"
)

type RegisterRequest struct {
	RoomID         uint32 `json:"room_id,omitempty"`
	BaseRouter     string `json:"baseRouter,omitempty"`
	DanmakuRouter  string `json:"danmakuRouter,omitempty"`
	GiftRouter     string `json:"giftRouter,omitempty"`
	GuardRouter    string `json:"guardRouter,omitempty"`
	MasterRouter   string `json:"masterRouter,omitempty"`
	AudienceRouter string `json:"audienceRouter,omitempty"`
	FansRouter     string `json:"fansRouter,omitempty"`
	CustomMessage  string `json:"customMessage,omitempty"`
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
				log.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.UserLevel, m.UserName, m.Message)
				callbackClient.SendMessage(callback.NewContext(m, callbackClient.Options.DanmakuRouter))
			case src := <-pool.Gift:
				g := message.NewGiftWithData(src)
				log.Printf("%s %s valued %d gift %s * %v\n", g.UserName, g.Action, g.Price, g.GiftName, g.Number)
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

		client := callback.NewClient(
			request.BaseRouter,
			request.DanmakuRouter,
			request.GiftRouter,
			request.GuardRouter,
			request.MasterRouter,
			request.AudienceRouter,
			request.FansRouter,
			request.CustomMessage,
		)
		messageChan := pool.NewPoolWithHandler(WebsocketHandler(client))

		wsClient := startWebsocketServer(request.RoomID, messageChan)

		ctx.JSON(200, RoomInfo{
			RoomID:     wsClient.RoomID,
			UpUID:      wsClient.UpUID,
			Title:      wsClient.Title,
			Online:     wsClient.Online,
			Tags:       wsClient.Tags,
			LiveStatus: wsClient.LiveStatus,
		})
	}
}

func startWebsocketServer(realRoomID uint32, messageChan *pool.Pool) info.Room {
	client := websocket.NewClientWithMessageChan(realRoomID, messageChan)
	GetDefaultChannelService().AddChannel(realRoomID, client)

	success := false
	for i := 0; i < 10; i++ {
		success = client.Serve()
		if success {
			break
		}
	}

	return client.GetRoomInfo()
}
