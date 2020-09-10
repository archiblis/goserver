package main

import (
	//"aqua"
	"time"
	//l4g "github.com/jeanphorn/log4go"
)

type GameLogic struct {
}

func NewGameLogic() *GameLogic {
	gl := &GameLogic{}
	return gl
}

func (this *GameLogic) Loop() {
	for {
		g_nowTime = time.Now()
		select {
		case <-g_sceneSessioner.closeChan:
			go g_sceneSessioner.Connect()
		case <-g_sceneSessioner.connectChan:
			g_sceneSessioner.Login()
		case msg := <-g_sceneSessioner.readMsgQueue:
			g_sceneSessioner.cmdM.Dispatcher(msg.gs, msg.ph, msg.buf)
		default:
			g_timerLoop.Loop(g_nowTime.Unix())
		}
	}
}
