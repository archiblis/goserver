package main

import (
	"aqua"
	//"command"
	//	"fmt"
	//"github.com/golang/protobuf/proto"
	l4g "github.com/jeanphorn/log4go"
	//"protos/cs"
)

type EMPTY struct{}

var (
	g_config     *aqua.Config
	g_gameServer *GameServer
	g_gameLogic  *GameLogic
	g_userMgr    *UserMgr
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

	g_gameServer = &GameServer{
		unlogin: make(map[*GameSessioner]uint64),
		users:   make(map[uint64]*GameSessioner),
		//gsUser:  make(map[*GameSessioner]uint64),
		cmdM:         aqua.NewCommandM(),
		messageQueue: make(chan *TcpMessage, 1024),
	}
	g_gameServer.InitCmd()

	g_gameLogic = NewGameLogic()
	g_userMgr = NewUserMgr()
}

func main() {
	defer l4g.Close()
	//l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	//l4g.AddFilter("file", l4g.INFO, l4g.NewFileLogWriter("server.log", true, true))
	flw := l4g.NewFileLogWriter("../log/server.log", false, false)
	//flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetFormat("[%d %t] [%L] (%S) %M")
	flw.SetRotate(true)
	flw.SetRotateSize(1 * 1024 * 1024) // 10M
	flw.SetRotateMaxBackup(10)         // 备份的日志文件数量
	flw.SetRotateDaily(true)           // 不按天回滚
	flw.SetSanitize(true)
	l4g.AddFilter("file", l4g.INFO, flw)

	l4g.Info("Game Init")
	GameInit()

	//go aqua.TCPServe(g_gameServer, g_config)
	go aqua.TCPServeReusePort(g_gameServer, g_config)
	g_gameLogic.Loop()
}
