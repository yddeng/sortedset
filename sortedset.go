package grank

import "github.com/yddeng/grank/skiplist"

// key type
type Key string

// skiplist.Interface
type Interface interface {
	skiplist.Interface
}

type element struct {
	key   Key
	value Interface
}

func (e *element) Less(other interface{}) bool {
	return e.value.Less(other.(*element).value)
}

type SortedSet struct {
	dict map[Key]*skiplist.Element
	zsl  *skiplist.SkipList
}

// New returns an initialized sortedset.
func New() *SortedSet { return new(SortedSet).Init() }

// Init initializes or clears sortedset z.
func (this *SortedSet) Init() *SortedSet {
	this.dict = map[Key]*skiplist.Element{}
	if this.zsl == nil {
		this.zsl = skiplist.New()
	} else {
		this.zsl.Init()
	}
	return this
}

// Len returns counts of elements
func (z *SortedSet) Len() int {
	return z.zsl.Len()
}

// Set is used to add or update an element
func (this *SortedSet) Set(key Key, v Interface) {
	if e, ok := this.dict[key]; ok {
		this.zsl.Remove(e)
	}
	ele := &element{
		key:   key,
		value: v,
	}
	e := this.zsl.Insert(ele)
	this.dict[key] = e
}

// Delete removes an element from the SortedSet
// by its key.
func (this *SortedSet) Delete(key Key) (ok bool) {
	if e, ok := this.dict[key]; ok {
		this.zsl.Remove(e)
		delete(this.dict, key)
		return true
	}
	return false
}

// GetData returns data stored in the map by its key
func (this *SortedSet) GetData(key Key) (data interface{}, ok bool) {
	if e, ok := this.dict[key]; ok {
		return e.Value().(*element).value, true
	}
	return nil, false
}

func (this *SortedSet) GetRank(key Key) int {
	if e, ok := this.dict[key]; ok {
		return e.Rank()
	}
	return 0
}

func (this *SortedSet) GetDataByRank(rank int) interface{} {
	if rank <= 0 || rank > len(this.dict) {
		return nil
	}
	e := this.zsl.GetElementByRank(rank)
	return e.Value().(*element).value
}

// Range implements ZRANGE
func (z *SortedSet) Range(start, end int, f func(key Key, value interface{})) {
	for e, span := z.zsl.GetElementByRank(start), end-start+1; span > 0 && e != nil; e, span = e.Next(), span-1 {
		elem := e.Value().(*element)
		f(elem.key, elem.value)
	}
}

// RevRange implements ZREVRANGE
func (z *SortedSet) RevRange(start, end int, f func(key Key, value interface{})) {
	for e, span := z.zsl.GetElementByRank(end), end-start+1; span > 0 && e != nil; e, span = e.Prev(), span-1 {
		elem := e.Value().(*element)
		f(elem.key, elem.value)
	}
}
