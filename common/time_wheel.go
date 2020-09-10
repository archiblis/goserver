package common

import (
	"time"

	l4g "github.com/jeanphorn/log4go"
)

func noneFunc() {
	l4g.Info("noneFunc")
}

// 时间轮
// 仅处理心跳包用

const (
	WHEEL_SIZE int = 60
)

/*
type WheelUser struct {
	id uint64
}

func NewWheelUser(id uint64) *WheelUser {
	nwu := &WheelUser{
		id: id,
	}
	return nwu
}
*/
type Wheel struct {
	ids map[uint64]EMPTY
}

func NewWheel() *Wheel {
	nw := &Wheel{
		ids: make(map[uint64]EMPTY),
	}
	return nw
}

type TimeWheel struct {
	size   int64
	wheels []*Wheel
	id_map map[uint64]int

	start_time int64
	cur_index  int
}

func NewTimeWheel(size int64) *TimeWheel {
	ntw := &TimeWheel{
		size:       size,
		wheels:     make([]*Wheel, size),
		id_map:     make(map[uint64]int),
		start_time: 0,
		cur_index:  0,
	}

	for i := 0; i < int(size); i++ {
		ntw.wheels[i] = NewWheel()
	}
	return ntw
}

func (this *TimeWheel) Start() {
	this.start_time = time.Now().Unix()
	this.cur_index = 0 // this.size - 1
}

func (this *TimeWheel) Run(now int64) (ids []uint64) {
	if this.start_time == now {
		return
	}
	if this.start_time > now {
		l4g.Error("你妈死了，往回改时间")
		return
	}

	move_steps := now - this.start_time
	if move_steps+int64(this.cur_index) >= this.size {
		this.start_time = now
	}

	if move_steps > this.size {
		move_steps = this.size
	}

	for i := 1; i <= int(move_steps); i++ {
		index := this.cur_index + i
		if index >= int(this.size) {
			index = 0
		}
		if len(this.wheels) == 0 {
			continue
		}
		for id, _ := range this.wheels[index].ids {
			ids = append(ids, id)
			delete(this.id_map, id)
		}
	}

	return
}

func (this *TimeWheel) Update(id uint64) {
	if index, suc := this.id_map[id]; suc {
		delete(this.id_map, id)
		delete(this.wheels[index].ids, id)
	}

	this.wheels[this.cur_index].ids[id] = Empty
	this.id_map[id] = this.cur_index
}
