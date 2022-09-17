package main

import (
	"fmt"
	"log"

	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/pool"
	"github.com/sunist-c/bililive-danmaku-backend/service/api"
	"github.com/sunist-c/bililive-danmaku-backend/websocket"
)

func main() {
	model.WaitGroup.Add(1)
	if model.GetGlobalConfig().BackendMode {
		api.StartHttpService()
	} else {
		log.Printf("Please input your room-id\n")
		var room uint32
		_, _ = fmt.Scanf("%d", &room)
		if room == 0 {
			log.Printf("Bad room-id, exit...\n")
			return
		} else {
			log.Printf("Try to connect to room %d\nYou can exit the program by type 'exit' any time\n", room)

			client := websocket.NewClientWithHandler(room, pool.EmptyHandler())
			success := false
			for i := 0; i < 10; i++ {
				success = client.Serve()
				if success {
					break
				}
			}

			go func() {
				var x string
				_, _ = fmt.Scanf("%v\n", &x)
				if x == "exit" || x == "exit\n" || x == "exit\r\n" {
					model.WaitGroup.Done()
				}
			}()
		}
	}
	model.WaitGroup.Wait()
}
