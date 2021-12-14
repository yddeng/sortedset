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
	sl   *skiplist.SkipList
}

// New returns an initialized sortedset.
func New() *SortedSet { return new(SortedSet).Init() }

// Init initializes or clears sortedset z.
func (z *SortedSet) Init() *SortedSet {
	z.dict = map[Key]*skiplist.Element{}
	z.sl = skiplist.New()
	return z
}

// Len returns counts of elements
func (z *SortedSet) Len() int {
	return z.sl.Len()
}

// Set is used to add or update an element,returns
// the rank where it would be inserted.
func (z *SortedSet) Set(key Key, v Interface) (rank int) {
	if e, ok := z.dict[key]; ok {
		z.sl.Remove(e)
	}
	var e *skiplist.Element
	ele := &element{key: key, value: v}
	e, rank = z.sl.Insert(ele)
	z.dict[key] = e
	return
}

// Delete removes an element from the SortedSet
// by its key.
func (z *SortedSet) Delete(key Key) (value interface{}) {
	if e, ok := z.dict[key]; ok {
		z.sl.Remove(e)
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

// Select returns key,value by rank.
func (z *SortedSet) Select(rank int) (key Key, value interface{}) {
	if rank <= 0 || rank > len(z.dict) {
		return
	}
	e := z.sl.Select(rank)
	key, value = e.Value().(*element).key, e.Value().(*element).value
	return
}

// Range implements ZRANGE
// If f returns false, range stops the iteration.
func (z *SortedSet) Range(start, end int, f func(rank int, key Key, value interface{}) bool) {
	var elem *element
	for e, i := z.sl.Select(start), start; i <= end && e != nil; e, i = e.Next(), i+1 {
		elem = e.Value().(*element)
		if !f(i, elem.key, elem.value) {
			return
		}
	}
}

// RevRange implements ZREVRANGE
// If f returns false, range stops the iteration.
func (z *SortedSet) RevRange(start, end int, f func(rank int, key Key, value interface{}) bool) {
	var elem *element
	for e, i := z.sl.Select(end), end; i >= start && e != nil; e, i = e.Prev(), i-1 {
		elem = e.Value().(*element)
		if !f(i, elem.key, elem.value) {
			return
		}
	}
}

// Search implements skiplist.Search
func (z *SortedSet) Search(n int, f func(i Interface) bool) (rank int) {
	rank = z.sl.Search(n, func(i skiplist.Interface) bool {
		return f(i.(*element).value)
	})
	return
}

// WouldBeInserted implements skiplist.WouldBeInserted
func (z *SortedSet) WouldBeInserted(v Interface) (rank int) {
	return z.sl.WouldBeInserted(&element{value: v})
}
