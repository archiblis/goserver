package main

import (
	"aqua"
	"command"
	"protos/cs"

	"github.com/golang/protobuf/proto"
	l4g "github.com/jeanphorn/log4go"
)

func RegisterGameServer() {
	l4g.Info("RegisterGameServer")
	g_gameServer.cmdM.Register(command.C2S_LOGIN_REQUEST, C2SLogin)
	g_gameServer.cmdM.Register(command.C2S_HEART_BEAT_REQUEST, C2SHeartBeat)
}

func C2SLogin(s aqua.Sessioner, ph *aqua.PackHead, buff []byte) bool {
	l4g.Info("C2SLogin")

	/*
		user := g_userMgr.GetUserBySessioner(&s)
		if user == nil {
			l4g.Error("not found user")
			return false
		}
	*/
	return true
}

func C2SHeartBeat(s aqua.Sessioner, ph *aqua.PackHead, buff []byte) bool {
	var request cs.C2SHeartBeat
	err := proto.Unmarshal(buff, &request)
	if err != nil {
		l4g.Error("C2SHeartBeat unmarshal error %v", err)
		return false
	}

	l4g.Info("C2SHeartBeat %s", request.String())
	msg := &cs.S2CHeartBeat{
		Index: request.Index,
	}

	gs := g_gameServer.getSessioner(ph.Uid)
	gs.Write(command.C2S_HEART_BEAT_RESPONE, msg)
	return true
}
