package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/bililive-danmaku-backend/service/system"

	"github.com/sunist-c/bililive-danmaku-backend/callback"
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
	UpdateRoom     string `json:"updateRoom,omitempty"`
	ExitRouter     string `json:"exitRouter,omitempty"`
}

func RegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &RegisterRequest{}
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		option := callback.ClientOptions{
			BaseRouter:     request.BaseRouter,
			DanmakuRouter:  request.DanmakuRouter,
			GiftRouter:     request.GiftRouter,
			GuardRouter:    request.GuardRouter,
			MasterRouter:   request.MasterRouter,
			AudienceRouter: request.AudienceRouter,
			FansRouter:     request.FansRouter,
			CustomMessage:  request.CustomMessage,
			UpdateRoom:     request.UpdateRoom,
			ExitRouter:     request.ExitRouter,
		}

		messageService := system.NewMessageService()
		messageService = messageService.InitializeMessageService(request.RoomID, &option, system.WebhookImplementation(messageService))

		channelService.AddChannel(request.RoomID, messageService)
	}
}
