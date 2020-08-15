package aqua

import (
	//	l4g "base/log4go"
	l4g "github.com/jeanphorn/log4go"
)

type PackHead struct {
	Length uint32
	Cmd    uint32
	Uid    uint64
	Sid    uint32
}

func DecodePackHead(buf []byte, ph *PackHead) bool {
	if len(buf) < 20 {
		l4g.Error("[PackHead] decode size no enough size: %v", len(buf))
		return false
	}
	ph.Length = DecodeUint32(buf[0:])
	ph.Cmd = DecodeUint32(buf[4:])
	ph.Uid = DecodeUint64(buf[8:])
	ph.Sid = DecodeUint32(buf[16:])
	return true
}

func EncodePackHead(buf []byte, ph *PackHead) bool {
	if len(buf) < 20 {
		l4g.Error("[PackHead] encode size no enough size: %v", len(buf))
		return false
	}
	EncodUint32(ph.Length, buf[0:])
	EncodUint32(ph.Cmd, buf[4:])
	EncodUint64(ph.Uid, buf[8:])
	EncodUint32(ph.Sid, buf[16:])
	return true
}
