package main

import (
	"aqua"
	"common"
	"math/rand"
	"time"

	l4g "github.com/jeanphorn/log4go"
)

type EMPTY = struct{}

var (
	g_config *aqua.Config

	// g_game_server *GameServer
	g_sceneSessioner *ServerSessioner
	g_gameLogic      *GameLogic
	g_timerLoop      *TimerLoop
	g_nowTime        time.Time
)

func GameInit() {
	g_config = &aqua.Config{
		Address:           "127.0.0.1:10001",
		MaxReadMsgSize:    65536,
		ReadMsgQueueSize:  0, //1024,
		ReadTimeOut:       60,
		MaxWriteMsgSize:   65536,
		WriteMsgQueueSize: 1024,
		WriteTimeOut:      60,
	}

	g_sceneSessioner = NewServerSessioner("scene")
	RegisterGameServer()
	g_sceneSessioner.closeChan <- struct{}{}
	g_gameLogic = NewGameLogic()
	g_timerLoop = NewTimerLoop()
	g_nowTime = time.Now()
}

func main() {
	defer l4g.Close()
	//l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	//l4g.AddFilter("file", l4g.INFO, l4g.NewFileLogWriter("client.log", true, true))

	flw := l4g.NewFileLogWriter("../log/client.log", false, false)
	//flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetFormat("[%d %t] [%L] (%S) %M")
	flw.SetRotate(true)
	flw.SetRotateSize(1 * 1024 * 1024) // 10M
	flw.SetRotateMaxBackup(10)         // 备份的日志文件数量
	flw.SetRotateDaily(true)           // 不按天回滚
	flw.SetSanitize(true)
	l4g.AddFilter("file", l4g.INFO, flw)

	//go Tcp()

	l4g.Info("Client Init")
	GameInit()

	if false {
		heap := common.NewHeap(common.Min)
		rand.Seed(time.Now().Unix())
		for i := 0; i < 15; i++ {
			r := int64(rand.Int31n(100))
			l4g.Info("rand %d", r)
			node := common.NewNode(r, nil, nil)
			heap.Insert(node)
		}

		heap.Output()
		for i := 0; i < 15; i++ {
			node := heap.Delete(3)
			if node == nil {
				break
			}
			l4g.Info("------------------")
			l4g.Info("delet root %d %d", node.Key(), node.Pos())
			heap.Output()
			l4g.Info("------------------")
		}
	}
	if false {
		node := common.NewNode(111, UserCallback, &UserTimerParam{"abc"})
		g_timerLoop.heap.Insert(node)
		g_timerLoop.heap.DeleteRoot()
	}

	g_gameLogic.Loop()

	for {
		time.Sleep(time.Second * 100)
	}
	l4g.Info("Client Connect")
}
