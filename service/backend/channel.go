package backend

import (
	"log"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/sunist-c/bililive-danmaku-backend/websocket"
)

var (
	channelService *ChannelService
)

type Channel struct {
	mu    *sync.RWMutex
	rooms map[uint32]*websocket.Client
}

type ChannelService struct {
	data *Channel
}

func (s *ChannelService) AddChannel(realRoomID uint32, client *websocket.Client) (success bool) {
	s.data.mu.Lock()
	defer s.data.mu.Unlock()

	if _, ok := s.data.rooms[realRoomID]; ok {
		log.Printf("add existed room %v\n", realRoomID)
		return false
	} else {
		s.data.rooms[realRoomID] = client
		return true
	}
}

func (s *ChannelService) RemoveChannel(realRoomID uint32) (success bool) {
	s.data.mu.Lock()
	defer s.data.mu.Unlock()

	if client, ok := s.data.rooms[realRoomID]; !ok {
		log.Printf("try to remove a unregisted room: %v\n", realRoomID)
		return true
	} else {
		client.Stop()
		delete(s.data.rooms, realRoomID)
		return true
	}
}

func GetDefaultChannelService() *ChannelService {
	if channelService == nil {
		channelService = &ChannelService{
			data: &Channel{
				mu:    &sync.RWMutex{},
				rooms: map[uint32]*websocket.Client{},
			},
		}
	}

	return channelService
}

func RemoveChannelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		channel, ok := ctx.Params.Get("channel_id")
		if !ok {
			ctx.AbortWithStatus(400)
		}

		roomID, err := strconv.ParseUint(channel, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(400)
		}

		go channelService.RemoveChannel(uint32(roomID))
		ctx.JSON(200, "success")
	}
}
