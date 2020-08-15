package aqua

import (
	//	l4g "log4go"
	l4g "github.com/jeanphorn/log4go"
)

type Services func(Sessioner, *PackHead, []byte) bool

type CommandM struct {
	cmdm map[uint32]Services
}

func NewCommandM() *CommandM {
	return &CommandM{
		cmdm: make(map[uint32]Services),
	}
}

func (this *CommandM) Register(id uint32, service Services) {
	this.cmdm[id] = service
}

func (this *CommandM) Dispatcher(session Sessioner, ph *PackHead, data []byte) bool {
	if cmd, exist := this.cmdm[ph.Cmd]; exist {
		return cmd(session, ph, data)
	}
	l4g.Error("[Command] no find cmd: %d %d", ph.Sid, ph.Cmd)
	return false
}
