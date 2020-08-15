package aqua

import (
	l4g "github.com/jeanphorn/log4go"
	"net"
)

type Server interface {
	NewSession() Sessioner
	Close()
}

// 启动 tcp server
func TCPServe(src Server, conf *Config) {
	defer src.Close()

	l, e := net.Listen("tcp", conf.Address)
	if e != nil {
		l4g.Error("[TCPServer] listen error: %v", e)
		panic(e.Error())
	}

	defer l.Close()

	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				continue
			}
			l4g.Error("[TCPServer] accept error: %v", e)
			return
		}
		newBroker(src.NewSession(), conf).serve(rw)
	}
}
