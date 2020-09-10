package main

import (
	"aqua"
	"command"

	"github.com/golang/protobuf/proto"
	l4g "github.com/jeanphorn/log4go"
	"protos/cs"
)

// TcpMessage
type TcpMessage struct {
	ph  *aqua.PackHead
	buf []byte
	gs  *GameSessioner
}

// GameServer
type GameServer struct {
	unlogin map[*GameSessioner]uint64
	users   map[uint64]*GameSessioner
	//gsUser       map[*GameSessioner]uint64
	cmdM         *aqua.CommandM
	messageQueue chan *TcpMessage
}

func (this *GameServer) NewSession() aqua.Sessioner {
	gs := &GameSessioner{server: this}
	this.unlogin[gs] = 0
	return gs
}

func (this *GameServer) InitCmd() {
	RegisterGameServer()
}

func (this *GameServer) Close() {
}

func (this *GameServer) getSessioner(id uint64) *GameSessioner {
	if v, suc := this.users[id]; suc {
		return v
	}
	return nil
}

func (this *GameServer) SessionLogin(id uint64, gs *GameSessioner) bool {
	if _, ok := this.unlogin[gs]; ok == false {
		l4g.Error("game server not find session, uid %ld", id)
		return false
	}

	this.unlogin[gs] = id
	this.users[id] = gs

	ph := &aqua.PackHead{
		Cmd: command.C2S_LOGIN_RESPONE,
		Uid: 1,
		Sid: 1,
	}
	login := &cs.S2CLogin{
		Ret:  1,
		Id:   1,
		Name: "Cpp",
	}
	gs.broker.Write(ph, login)

	l4g.Info("session login %d", id)
	// todo check user
	return true
}

func (this *GameServer) Login(id uint64, gs *GameSessioner) bool {
	if gs == nil {
		l4g.Error("game server login gs nil")
		return false
	}
	if _, ok := this.unlogin[gs]; !ok {
		l4g.Error("game server unlogin error")
		return false
	}

	delete(this.unlogin, gs)
	this.users[id] = gs
	//this.gsUser[gs] = id
	g_userMgr.AddUser(id, gs)

	l4g.Info("login %d", id)
	return true
}

// GameSessioner
type GameSessioner struct {
	broker *aqua.Broker
	server *GameServer
}

func (this *GameSessioner) Init(broker *aqua.Broker) {
	this.broker = broker
}

func (this *GameSessioner) Process(ph *aqua.PackHead, buf []byte) {
	l4g.Info("Process cmd %d uid %d", ph.Cmd, ph.Uid)
	if ph.Cmd == command.C2S_LOGIN_REQUEST {
		g_gameServer.SessionLogin(ph.Uid, this)
		return
	}

	this.server.messageQueue <- &TcpMessage{ph, buf, this}
}

func (this *GameSessioner) Close() {

}

func (this *GameSessioner) ReadLimt() {
	// todo
}

func (this *GameSessioner) Write(cmd uint32, message interface{}) {
	ph := &aqua.PackHead{}
	ph.Cmd = cmd
	ph.Uid = 1
	ph.Sid = 1

	msg, suc := message.(proto.Message)
	if suc {
		this.broker.Write(ph, msg)
	} else {
		l4g.Error("write msg error %d", cmd)
	}
}
