package aqua

import (
	//	l4g "base/log4go"
	l4g "github.com/jeanphorn/log4go"
)

type Config struct {
	Address string // ip:prot
	// read config
	MaxReadMsgSize   int
	ReadMsgQueueSize int
	ReadTimeOut      int
	// write config
	MaxWriteMsgSize   int
	WriteMsgQueueSize int
	WriteTimeOut      int
}

func (this *Config) Check() bool {
	if this.MaxReadMsgSize == 0 {
		l4g.Error("[Config] MaxWriteMsgSize error")
		return false
	}
	if this.WriteMsgQueueSize == 0 {
		l4g.Error("[Config] WriteMsgQueueSize error")
		return false
	}
	if this.MaxReadMsgSize == 0 {
		l4g.Error("[Config] MaxReadMsgSize error")
		return false
	}
	return true
}
