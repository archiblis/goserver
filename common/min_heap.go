package common

import (
	//"fmt"
	l4g "github.com/jeanphorn/log4go"
)

//
// true : key1 需在堆的上方
// false: key1 需要堆的下方
type HeapCompare func(key1 int64, key2 int64) bool

// callback
type HeapCallBack func(param interface{})

func Min(key1 int64, key2 int64) bool {
	return key1 < key2
}

func Max(key1 int64, key2 int64) bool {
	return key1 > key2
}

type Node struct {
	key      int64
	id       int64
	pos      int
	callback HeapCallBack
	param    interface{}
}

func (this *Node) Key() int64 {
	return this.key
}

func (this *Node) Id() int64 {
	return this.id
}

func (this *Node) Pos() int {
	return this.pos
}

func (this *Node) Callback() {
	if this.callback != nil {
		this.callback(this.param)
	}
}

func NewNode(key int64, callback HeapCallBack, param interface{}) *Node {
	return &Node{key, 0, 0, callback, param}
}

type Heap struct {
	comp   HeapCompare
	nodes  []*Node
	len    int
	size   int
	id_map map[int64]*Node
	id     int64
}

func NewHeap(comp HeapCompare) *Heap {
	len := 1
	heap := &Heap{
		comp:   comp,
		nodes:  make([]*Node, len),
		len:    len,
		size:   0,
		id_map: make(map[int64]*Node),
		id:     1,
	}
	return heap
}

func (this *Heap) Len() int {
	return this.len
}

func (this *Heap) Size() int {
	return this.size
}

func (this *Heap) Comp(key1 int64, key2 int64) bool {
	return this.comp(key1, key2)
}

func (this *Heap) Parent(pos int) int {
	if pos <= 0 {
		return 0
	}
	return (pos - 1) / 2
}

func (this *Heap) Left(pos int) int {
	return (pos+1)*2 - 1
}

func (this *Heap) Right(pos int) int {
	return (pos + 1) * 2
}

func (this *Heap) Resize() int {
	nodes := make([]*Node, this.len)
	this.nodes = append(this.nodes, nodes...)
	this.len *= 2
	return this.len
}

func (this *Heap) Swap(pos1, pos2 int) {
	node := this.nodes[pos1]
	this.nodes[pos1] = this.nodes[pos2]
	this.nodes[pos2] = node
	this.nodes[pos1].pos = pos1
	this.nodes[pos2].pos = pos2
}

func (this *Heap) Insert(node *Node) {
	node.pos = this.size
	cur_pos := this.size
	this.size++
	node.id = this.id
	this.id++

	this.id_map[node.id] = node
	this.nodes[cur_pos] = node
	this.Up(cur_pos)

	if this.size == this.len {
		this.Resize()
	}
}

func (this *Heap) Up(pos int) {
	if this.size == 0 || pos < 0 || pos >= this.size {
		return
	}

	node := this.nodes[pos]
	root := this.Parent(pos)
	for root != pos {
		if this.comp(node.key, this.nodes[root].key) {
			this.Swap(pos, root)
			pos = root
			node = this.nodes[pos]
			root = this.Parent(pos)
		} else {
			break
		}
	}
}

func (this *Heap) DeleteRoot() *Node {
	node := this.Delete(0)
	if node != nil {
		node.Callback()
	}

	return node
}

func (this *Heap) Delete(pos int) *Node {
	if this.size == 0 || pos < 0 || pos >= this.size {
		return nil
	}

	node := this.nodes[pos]
	delete(this.id_map, node.id)
	this.Swap(pos, this.size-1)
	this.size--
	this.Up(pos)
	p := this.Down(pos)
	if false {
		l4g.Info("delete pos %d", p)
	}
	return node
}

func (this *Heap) DeleteById(id int64) *Node {
	if v, suc := this.id_map[id]; suc == true {
		return this.Delete(v.pos)
	}
	return nil
}

func (this *Heap) Down(pos int) int {
	p := -1
	for {
		l := this.Left(pos)
		r := this.Right(pos)
		if l >= this.size {
			break
		} else if r >= this.size {
			p = l
		} else if this.comp(this.nodes[r].key, this.nodes[l].key) {
			p = r
		} else {
			p = l
		}

		if this.comp(this.nodes[p].key, this.nodes[pos].key) {
			this.Swap(p, pos)
			pos = p
			continue
		} else {
			break
		}
	}

	return p
}

func (this *Heap) GetRoot() *Node {
	if this.nodes[0] != nil {
		return this.nodes[0]
	}

	return nil
}

func (this *Heap) Output() {
	start := 0
	end := this.Left(start)
	l4g.Info("len %d size %d", this.len, this.size)
	for {
		if this.size <= end {
			s := this.nodes[start:this.size]
			for _, v := range s {
				l4g.Info("%d %d", v.key, v.pos)
			}
			l4g.Info("------------------")
			break
		}
		s := this.nodes[start:end]
		for _, v := range s {
			l4g.Info("%d %d", v.key, v.pos)
		}

		l4g.Info("------------------")
		start = end
		end = this.Left(start)
	}
}
