package main

import (
	"command"
	"common"
	"protos/cs"
	"time"

	l4g "github.com/jeanphorn/log4go"
)

const (
	HEART_BEAT_SECOND int64 = 60
)

//type TimerServices func()

type Timer interface {
	Timeout(int64)
	IsLoop() bool
}

type TimerLoop struct {
	heap *common.Heap
	tm   time.Time
}

func NewTimerLoop() *TimerLoop {
	tl := &TimerLoop{
		heap: common.NewHeap(common.Min),
	}
	return tl
}

func (this *TimerLoop) Loop(now int64) {
	for {
		node := this.heap.GetRoot()
		if node == nil || now < node.Key() {
			break
		}
		//l4g.Info("timer loop %d %d", now, node.Key())
		this.heap.DeleteRoot()
	}
}

func UserCallback(param interface{}) {
	p := param.(*UserTimerParam)
	l4g.Info("aaaa %s", p.str)
}

type UserTimerParam struct {
	str string
}

func HeartBeatTimerCallBack(param interface{}) {
	l4g.Info("HeartBeatTimerCallBack %d", g_nowTime.Unix())
	if p, suc := param.(*HeartBeatParam); suc {
		index := p.index
		msg := &cs.C2SHeartBeat{
			Index: index,
		}

		g_sceneSessioner.Write(command.C2S_HEART_BEAT_REQUEST, msg)

		index += 1
		node := common.NewNode(g_nowTime.Unix()+HEART_BEAT_SECOND, HeartBeatTimerCallBack, &HeartBeatParam{index})
		g_timerLoop.heap.Insert(node)
	} else {
		l4g.Error("heart beat timer error")
	}
}

type HeartBeatParam struct {
	index int32
}
