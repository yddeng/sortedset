package sortedset

import (
	"github.com/yddeng/sortedset/skiplist"
)

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
func (z *SortedSet) Init() *SortedSet {
	z.dict = map[Key]*skiplist.Element{}
	if z.zsl == nil {
		z.zsl = skiplist.New()
	} else {
		z.zsl.Init()
	}
	return z
}

// Len returns counts of elements
func (z *SortedSet) Len() int {
	return z.zsl.Len()
}

// Set is used to add or update an element
func (z *SortedSet) Set(key Key, v Interface) {
	if e, ok := z.dict[key]; ok {
		z.zsl.Remove(e)
	}
	ele := &element{
		key:   key,
		value: v,
	}
	e := z.zsl.Insert(ele)
	z.dict[key] = e
}

// Delete removes an element from the SortedSet
// by its key.
func (z *SortedSet) Delete(key Key) (value interface{}) {
	if e, ok := z.dict[key]; ok {
		z.zsl.Remove(e)
		delete(z.dict, key)
		return e.Value().(*element).value
	}
	return nil
}

// GetValue returns value in the map by its key
func (z *SortedSet) GetValue(key Key) (value interface{}, ok bool) {
	if e, ok := z.dict[key]; ok {
		return e.Value().(*element).value, true
	}
	return nil, false
}

// GetRank returns the rank of the element specified by key
func (z *SortedSet) GetRank(key Key) int {
	if e, ok := z.dict[key]; ok {
		return e.Rank()
	}
	return 0
}

// GetByRank returns key,value by rank.
func (z *SortedSet) GetByRank(rank int) (Key, interface{}) {
	if rank <= 0 || rank > len(z.dict) {
		return "", nil
	}
	e := z.zsl.GetElementByRank(rank)
	elem := e.Value().(*element)
	return elem.key, elem.value
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
