package pool

import (
	"log"
	"strings"

	"github.com/sunist-c/bililive-danmaku/model"
	"github.com/sunist-c/bililive-danmaku/model/message"
)

func EmptyHandler() func(pool *Pool) {
	return func(pool *Pool) {
		for {
			select {
			case uc := <-pool.Unknown:
				if cmd := model.Json.Get(uc, "cmd").ToString(); message.Command(cmd) == message.CommandRoomFocusedChange {
					fans := model.Json.Get(uc, "data", "fans").ToInt()
					log.Printf("room fans changed: %v\n", fans)
				}
			case src := <-pool.Danmaku:
				m := message.NewDanmakuWithData(src)
				log.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.UserLevel, m.UserName, m.Message)
			case src := <-pool.Gift:
				g := message.NewGiftWithData(src)
				log.Printf("%s %s valued %d gift %s * %v\n", g.UserName, g.Action, g.Price, g.GiftName, g.Number)
			case src := <-pool.Audience:
				name := model.Json.Get(src, "data", "uname").ToString()
				log.Printf("welcome master %s entered room\n", name)
			case src := <-pool.Guard:
				name := model.Json.Get(src, "data", "username").ToString()
				log.Printf("welcome guard %s entered room\n", name)
			case src := <-pool.Master:
				cw := model.Json.Get(src, "data", "copy_writing").ToString()
				cw = strings.Replace(cw, "<%", "", 1)
				cw = strings.Replace(cw, "%>", "", 1)
				log.Printf("%s\n", cw)
			}
		}
	}
}
