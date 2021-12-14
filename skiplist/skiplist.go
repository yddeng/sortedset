package skiplist

import (
	"math/rand"
)

const SKIPLIST_MAXLEVEL = 32
const SKIPLIST_BRANCH = 4

func randomLevel() int {
	level := 1
	for level < SKIPLIST_MAXLEVEL && (rand.Int31()&0xFFFF)%SKIPLIST_BRANCH == 0 {
		level += 1
	}
	return level
}

type Interface interface {
	Less(other interface{}) bool
}

type link struct {
	prev *Element
	next *Element
	skip int // 跳多少到下一节点
}

type Element struct {
	value Interface
	links []link
	sl    *SkipList
}

func (e *Element) Value() Interface {
	if e == nil {
		return nil
	}
	return e.value
}

// Next returns the next skiplist element or nil.
func (e *Element) Next() *Element {
	if e == nil || e.sl == nil || e.links[0].next == e.sl.tail {
		return nil
	}
	return e.links[0].next
}

// Prev returns the previous skiplist element of nil.
func (e *Element) Prev() *Element {
	if e == nil || e.sl == nil || e.links[0].prev == e.sl.head {
		return nil
	}
	return e.links[0].prev
}

func (e *Element) Rank() int {
	if e == nil || e.sl == nil {
		return 0
	}
	return e.sl.GetRank(e)
}

// newElement returns an initialized element.
func newElement(sl *SkipList, level int, v Interface) *Element {
	return &Element{
		value: v,
		links: make([]link, level),
		sl:    sl,
	}
}

type SkipList struct {
	head   *Element
	tail   *Element
	update []*Element
	rank   []int
	len    int
	level  int
}

// New returns an initialized skiplist.
func New() *SkipList {
	sl := new(SkipList)
	sl.level = 1
	sl.head = newElement(sl, SKIPLIST_MAXLEVEL, nil)
	sl.tail = newElement(sl, SKIPLIST_MAXLEVEL, nil)
	sl.update = make([]*Element, SKIPLIST_MAXLEVEL)
	sl.rank = make([]int, SKIPLIST_MAXLEVEL)

	for i := 0; i < SKIPLIST_MAXLEVEL; i++ {
		sl.head.links[i].next = sl.tail
		sl.tail.links[i].prev = sl.head
	}
	return sl
}

// Len returns the number of elements of skiplist sl.
// The complexity is O(1).
func (sl *SkipList) Len() int { return sl.len }

// Front returns the first element of skiplist sl or nil if the skiplist is empty.
func (sl *SkipList) Front() *Element {
	if sl.len == 0 {
		return nil
	}
	return sl.head.links[0].next
}

// Back returns the last element of skiplist sl or nil if the skiplist is empty.
func (sl *SkipList) Back() *Element {
	if sl.len == 0 {
		return nil
	}
	return sl.tail.links[0].prev
}

// Insert inserts v, increments sl.length, and returns a new element of wrap v
// and the rank where it would be inserted.
func (sl *SkipList) Insert(v Interface) (x *Element, rank int) {
	x = sl.head
	rank = 1
	sl.rank[sl.level-1] = 0
	for i := sl.level - 1; i >= 0; i-- {
		if i != sl.level-1 {
			sl.rank[i] = sl.rank[i+1]
		}
		for x.links[i].next != sl.tail && x.links[i].next.value.Less(v) {
			rank += x.links[i].skip
			sl.rank[i] += x.links[i].skip
			x = x.links[i].next
		}
		sl.update[i] = x
	}

	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			sl.update[i] = sl.head
			sl.rank[i] = 0
			sl.head.links[i].skip = sl.len
		}
		sl.level = level
	}

	x = newElement(sl, level, v)
	for i := 0; i < level; i++ {
		x.links[i].prev = sl.update[i]
		x.links[i].next = sl.update[i].links[i].next
		x.links[i].next.links[i].prev = x
		x.links[i].prev.links[i].next = x

		x.links[i].skip = sl.update[i].links[i].skip + sl.rank[i] - sl.rank[0]
		sl.update[i].links[i].skip = sl.rank[0] - sl.rank[i] + 1
	}

	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		sl.update[i].links[i].skip++
	}

	sl.len += 1

	return
}

// Remove removes e from sl if e is an element of skiplist sl.
// It returns the element value e.value.
// The element must not be nil.
func (sl *SkipList) Remove(e *Element) interface{} {
	if e.sl == sl {
		// if e.sl == sl, sl must have been initialized when e was inserted
		// in sl or sl == nil (e is a zero Element) and sl.remove will crash

		x := e.links[len(e.links)-1].prev
		lv := len(e.links) - 1
		for i := 0; i < len(e.links); i++ {
			e.links[i].prev.links[i].skip += e.links[i].skip - 1
			e.links[i].next.links[i].prev = e.links[i].prev
			e.links[i].prev.links[i].next = e.links[i].next
			e.links[i].next = nil
			e.links[i].prev = nil
		}

		for lv < sl.level-1 {
			for ; x != sl.head && lv == len(x.links)-1; x = x.links[lv].prev {
			}

			for i := lv + 1; i < len(x.links) && lv < sl.level; i++ {
				x.links[i].skip -= 1
			}

			lv = len(x.links) - 1
			x = x.links[lv].prev
		}

		for sl.level > 1 && sl.head.links[sl.level-1].next == sl.tail {
			sl.level--
		}
		sl.len--
		e.sl = nil
	}
	return e.value
}

// GetRank return the position if e is an element of skiplist sl.
func (sl *SkipList) GetRank(e *Element) int {
	if e.sl != sl {
		return 0
	}
	/*
		该方法查找含有相同分数时，往往找到的是第一个或最后一个（Less是否取等）
			x := sl.head
			rank := 0
			for i := sl.level - 1; i >= 0; i-- {
				for x.links[i].next != sl.tail && x.links[i].next.value.Less(e.value) {
					rank += x.links[i].skip
					x = x.links[i].next
				}
			}
	*/

	rank := 0
	for lv, x := len(e.links)-1, e; x != sl.head; {
		x = x.links[lv].prev
		rank += x.links[lv].skip
		lv = len(x.links) - 1
	}
	return rank
}

// Select an element by ites rank. The rank argument needs bo be 1-based.
func (sl *SkipList) Select(rank int) *Element {
	if rank <= 0 || rank > sl.len {
		return nil
	}

	x := sl.head
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.links[i].next != sl.tail && traversed+x.links[i].skip <= rank {
			traversed += x.links[i].skip
			x = x.links[i].next
		}
		if traversed == rank {
			return x
		}
	}

	return nil
}

// Search calls f(i) only for rank in the range [1, n+1].
func (sl *SkipList) Search(n int, f func(i Interface) bool) (rank int) {
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.links[i].next != sl.tail && rank+x.links[i].skip <= n && f(x.links[i].next.value) {
			rank += x.links[i].skip
			x = x.links[i].next
		}
	}
	return rank + 1
}

// WouldBeInserted returns the rank where it would be inserted.
func (sl *SkipList) WouldBeInserted(v Interface) int {
	return sl.Search(sl.Len(), func(i Interface) bool {
		return i.Less(v)
	})
}
