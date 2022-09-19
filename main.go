package main

import (
	"fmt"
	"github.com/sunist-c/bililive-danmaku-backend/callback"
	"github.com/sunist-c/bililive-danmaku-backend/handler"
	"github.com/sunist-c/bililive-danmaku-backend/service/system"
	"log"

	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/service/api"
)

func main() {
	model.WaitGroup.Add(1)
	if model.GetGlobalConfig().BackendMode {
		api.StartHttpService()
	} else {
		c := callback.NewClient()
		c.Serve()
		c.Close()
		log.Printf("Please input your room-id\n")
		var room uint32
		_, _ = fmt.Scanf("%d", &room)
		if room == 0 {
			log.Printf("Bad room-id, exit...\n")
			return
		} else {
			log.Printf("Try to connect to room %d\nYou can exit the program by type 'exit' any time\n", room)

			messageService := system.NewMessageService()
			messageService.InitializeMessageService(room, &callback.ClientOptions{}, handler.DisplayImplementation())
			messageService.Serve()
		}
	}
	model.WaitGroup.Wait()
}
