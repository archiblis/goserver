package main

import (
	l4g "github.com/jeanphorn/log4go"
)

type GameLogic struct {
}

func NewGameLogic() *GameLogic {
	gl := &GameLogic{}
	return gl
}

func (this *GameLogic) Loop() {
	for {
		select {
		case msg := <-g_gameServer.messageQueue:
			l4g.Info("msg %v", msg.ph)
			g_gameServer.cmdM.Dispatcher(msg.gs, msg.ph, msg.buf)
		}
	}
}
