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

	zs.Set("hello", Score(2.2))
	zs.Set("world", Score(5.5))
	zs.Set(",", Score(1.1))
	zs.Set("how", Score(3.3))
	zs.Set("are", Score(3.3))
	zs.Set("you", Score(5.5))

	t.Log(zs.Len())
	rank := 1
	zs.Range(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank++
	})

	// revRange
	fmt.Println()
	rank = zs.Len()
	zs.RevRange(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank--
	})

	t.Log(zs.GetRank("hello"))
	t.Log(zs.GetValue("hello"))
	t.Log(zs.GetByRank(5))

	// update
	fmt.Println()
	zs.Set("hello", Score(6.6))
	rank = 1
	zs.Range(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank++
	})

	// delete
	fmt.Println()
	zs.Delete("hello")
	rank = 1
	zs.Range(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank++
	})

}
