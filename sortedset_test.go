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
	t.Log(zs.Set("hello", Score(6.6)))
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

	// search
	score := Score(3.3)
	t.Log(zs.Search(zs.Len(), func(i Interface) bool {
		return i.(Score) > score
	}))

	// would be inserted
	t.Log(zs.WouldBeInserted(score))
}

type User struct {
	name  string
	level int
	score int
}

/*
	分数从大到小排序
	分数相同，按照等级排序
	分数、等级都相同，按照名字排序
*/
func (this *User) Less(other interface{}) bool {
	o := other.(*User)
	if this.score > o.score {
		return true
	} else if this.score == o.score && this.level > o.level {
		return true
	} else if this.score == o.score && this.level == o.level && this.name > o.name {
		return true
	}
	return false
}

func TestNew2(t *testing.T) {
	zs := New()

	zs.Set("u1", &User{name: "u1", level: 2, score: 30})
	zs.Set("u2", &User{name: "u2", level: 2, score: 40})
	zs.Set("u3", &User{name: "u3", level: 3, score: 30})
	zs.Set("u4", &User{name: "u4", level: 3, score: 30})

	zs.Range(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(key, value.(*User))
	})
}
