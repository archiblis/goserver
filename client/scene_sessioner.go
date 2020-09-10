package main

import (
	"aqua"
	"command"
	"protos/cs"
	"time"
	//"sync/atomic"

	"github.com/golang/protobuf/proto"
	l4g "github.com/jeanphorn/log4go"
)

// TcpMessage
type TcpMessage struct {
	ph  *aqua.PackHead
	buf []byte
	gs  *GameSessioner
}

type GameSessioner struct {
	broker *aqua.Broker
	ss     *ServerSessioner
}

func NewGameSessioner(ss *ServerSessioner) *GameSessioner {
	s := &GameSessioner{
		ss: ss,
	}
	return s
}

func (this *GameSessioner) Init(broker *aqua.Broker) {
	l4g.Info("GameSessioner init")
	this.broker = broker
}

func (this *GameSessioner) Process(ph *aqua.PackHead, buf []byte) {
	l4g.Info("Process")
	this.ss.readMsgQueue <- &TcpMessage{ph, buf, this}
	/*
		var response cs.S2CLogin
		err := proto.Unmarshal(buf, &response)
		if err != nil {
			l4g.Error("proto error %v", err)
			return
		}
		l4g.Info("response %s", response.String())
	*/
}

func (this *GameSessioner) Close() {
	// todo
}

func (this *GameSessioner) ReadLimt() {
	// todo
}

// scene server mgr
type ServerSessioner struct {
	sess         *GameSessioner
	name         string
	init         bool
	closeChan    chan struct{}
	connectChan  chan struct{}
	readMsgQueue chan *TcpMessage
	cmdM         *aqua.CommandM
}

func NewServerSessioner(name string) *ServerSessioner {
	ss := &ServerSessioner{
		sess:         nil, //NewGameSessioner(ss),
		name:         name,
		init:         false,
		closeChan:    make(chan struct{}, 1),
		connectChan:  make(chan struct{}, 1),
		readMsgQueue: make(chan *TcpMessage, 1024),
		cmdM:         aqua.NewCommandM(),
	}
	ss.sess = NewGameSessioner(ss)
	return ss
}

func (this *ServerSessioner) Connect() {
	for {
		//if this.sess.broker == nil ||
		//	this.sess.broker.State() != aqua.StateConnected {
		suc := aqua.TCPClientServe(this.sess, g_config, 1*time.Second)
		if suc == false {
			time.Sleep(time.Second * 3)
			l4g.Info("try to connect %v", this.name)
		} else {
			l4g.Info("connect %v success", this.name)
			this.connectChan <- struct{}{}
			break
		}
		//}
	}
}

func (this *ServerSessioner) Login() {
	l4g.Info("Login")
	if false {
		ph := &aqua.PackHead{}
		ph.Cmd = command.C2S_LOGIN_REQUEST
		ph.Uid = 1
		ph.Sid = 1

		msg := &cs.C2SLogin{
			Id:   10,
			Name: "cpp",
		}

		this.sess.broker.Write(ph, msg)
		return
	}

	msg := &cs.C2SLogin{
		Id:   10,
		Name: "cpp",
	}

	this.Write(command.C2S_LOGIN_REQUEST, msg)
	//this.sess.broker.Write(ph, login)
}

func (this *ServerSessioner) Write(cmd uint32, message interface{}) {
	ph := &aqua.PackHead{}
	ph.Cmd = cmd
	ph.Uid = 1
	ph.Sid = 1

	msg, suc := message.(proto.Message)
	if suc {
		this.sess.broker.Write(ph, msg)
	} else {
		l4g.Error("write msg error %d", cmd)
	}
}
