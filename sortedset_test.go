package sortedset

import (
	"fmt"
	"testing"
)

type Score float64

func (this Score) Less(other interface{}) bool {
	return this >= other.(Score)
}

func TestNew(t *testing.T) {
	zs := New()
	// add
	zs.Set("hello", Score(2.2))
	zs.Set("world", Score(5.5))
	zs.Set(",", Score(1.1))
	zs.Set("how", Score(3.3))
	zs.Set("are", Score(3.3))
	zs.Set("you", Score(5.5))

	t.Log(zs.Len())
	// get rank by key
	t.Log(zs.GetRank("hello"))
	// get value by key
	t.Log(zs.GetValue("hello"))
	// get key,value by rank
	t.Log(zs.GetByRank(zs.Len()))

	// update
	fmt.Println()
	zs.Set("hello", Score(6.6))
	// range
	rank := 1
	zs.Range(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank++
	})

	// delete
	fmt.Println()
	t.Log(zs.Delete("hello"))
	// RevRange
	rank = zs.Len()
	zs.RevRange(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank--
	})
}
