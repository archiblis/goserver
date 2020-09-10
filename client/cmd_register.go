package main

import (
	"aqua"
	"command"
	"common"
	"protos/cs"

	"github.com/golang/protobuf/proto"
	l4g "github.com/jeanphorn/log4go"
)

func RegisterGameServer() {
	l4g.Info("RegisterGameServer")
	g_sceneSessioner.cmdM.Register(command.C2S_LOGIN_RESPONE, C2SLogin)
}

func C2SLogin(s aqua.Sessioner, ph *aqua.PackHead, buff []byte) bool {
	var response cs.S2CLogin
	err := proto.Unmarshal(buff, &response)
	if err != nil {
		l4g.Error("C2SLogin unmarshal error %v", err)
		return false
	}

	l4g.Info("C2SLogin %s", response.String())

	index := int32(1)
	node := common.NewNode(g_nowTime.Unix()+5, HeartBeatTimerCallBack, &HeartBeatParam{index})
	g_timerLoop.heap.Insert(node)
	return true
}
