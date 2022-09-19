package handler

import (
	"github.com/sunist-c/bililive-danmaku-backend/common/logging"
	"github.com/sunist-c/bililive-danmaku-backend/model"
	"github.com/sunist-c/bililive-danmaku-backend/model/message"
	"strings"
)

func DisplayImplementation() func(pool *Pool, exit chan struct{}) {
	return func(pool *Pool, exit chan struct{}) {
		for {
			select {
			case <-exit:
				return
			case uc := <-pool.Unknown:
				if cmd := model.Json.Get(uc, "cmd").ToString(); message.Command(cmd) == message.CommandRoomFocusedChange {
					fans := model.Json.Get(uc, "data", "fans").ToInt()
					logging.Info("room fans changed: %v", fans)
				}
			case src := <-pool.Danmaku:
				m := message.NewDanmakuWithData(src)
				logging.Info("Lv%d %s - Lv%d %s: %s", m.MedalLevel, m.MedalName, m.UserLevel, m.UserName, m.Message)
			case src := <-pool.Gift:
				g := message.NewGiftWithData(src)
				logging.Info("%s %s gift %s*%v, total valued %v", g.UserName, g.Action, g.GiftName, g.Number, g.Price)
			case src := <-pool.Audience:
				name := model.Json.Get(src, "data", "uname").ToString()
				logging.Info("welcome master %s entered room", name)
			case src := <-pool.Guard:
				name := model.Json.Get(src, "data", "username").ToString()
				logging.Info("welcome guard %s entered room", name)
			case src := <-pool.Master:
				cw := model.Json.Get(src, "data", "copy_writing").ToString()
				cw = strings.Replace(cw, "<%", "", 1)
				cw = strings.Replace(cw, "%>", "", 1)
				logging.Info("%s", cw)
			}
		}
	}
}
