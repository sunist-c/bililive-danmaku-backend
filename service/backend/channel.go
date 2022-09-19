package backend

import (
	"log"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/sunist-c/bililive-danmaku-backend/service/system"
)

var (
	channelService *ChannelService
)

type Channel struct {
	mu    *sync.RWMutex
	rooms map[uint32]*system.MessageService
}

type ChannelService struct {
	data *Channel
}

func (s *ChannelService) AddChannel(roomID uint32, service *system.MessageService) (success bool) {
	s.data.mu.Lock()
	defer s.data.mu.Unlock()

	if _, ok := s.data.rooms[roomID]; ok {
		log.Printf("add existed room %v\n", roomID)
		return false
	} else {
		s.data.rooms[roomID] = service
		return true
	}
}

func (s *ChannelService) RemoveChannel(roomID uint32) (success bool) {
	s.data.mu.Lock()
	defer s.data.mu.Unlock()

	if messageService, ok := s.data.rooms[roomID]; !ok {
		log.Printf("try to remove a unregisted room: %v\n", roomID)
		return true
	} else {
		delete(s.data.rooms, roomID)
		return messageService.Close()
	}
}

func (s *ChannelService) RemoveAllChannel() (success, total uint) {
	total, success = 0, 0
	for _, messageService := range s.data.rooms {
		total += 1
		if messageService.Close() {
			success += 1
		}
	}

	return success, total
}

func GetDefaultChannelService() *ChannelService {
	if channelService == nil {
		channelService = &ChannelService{
			data: &Channel{
				mu:    &sync.RWMutex{},
				rooms: map[uint32]*system.MessageService{},
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
