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
	t.Log(zs.Select(zs.Len()))

	// update
	fmt.Println()
	t.Log(zs.Set("hello", Score(6.6)))
	// range
	zs.Range(1, zs.Len(), func(rank int, key Key, value interface{}) bool {
		t.Log(rank, " -- ", key, value.(Score))
		return true
	})

	// delete
	fmt.Println()
	t.Log(zs.Delete("hello"))
	// RevRange
	zs.RevRange(1, zs.Len(), func(rank int, key Key, value interface{}) bool {
		t.Log(rank, " -- ", key, value.(Score))
		return true
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

	zs.Range(1, zs.Len(), func(rank int, key Key, value interface{}) bool {
		t.Log(rank, " -- ", key, value.(*User))
		return true
	})
}

type Int32 int32

// 递增序列，值相等按照先后顺序排列
func (i1 Int32) Less(other interface{}) bool {
	return i1 <= other.(Int32)
}

func TestGetByScore(t *testing.T) {
	l := New()
	l.Set("1", Int32(1))
	l.Set("2", Int32(2))
	l.Set("3", Int32(4))
	l.Set("4", Int32(5))
	l.Set("5", Int32(2))
	l.Set("6", Int32(3))
	l.Set("7", Int32(7))

	l.Range(1, l.Len(), func(rank int, key Key, value interface{}) bool {
		t.Log(rank, "--", key, value)
		return true
	})

	// 查找分数区间 [2,5]
	left := l.Search(l.Len(), func(i Interface) bool {
		return i.(Int32) < Int32(2)
	})
	right := l.Search(l.Len(), func(i Interface) bool {
		return i.(Int32) <= Int32(5)
	})

	right = right - 1 // search 是模拟插入动作，根据自定义的规则返回将要插入的位置，故减一

	if right > l.Len() {
		// 根据自定义的规则返回将要插入的位置，可能大于长度即最末插入
		right = l.Len()
	}
	t.Log(left, right)
	l.Range(left, right, func(rank int, key Key, value interface{}) bool {
		t.Log(rank, "--", key, value)
		return true
	})
}

// 获取用户排名附近用户
func TestGetAround(t *testing.T) {
	l := New()
	l.Set("1", Int32(1))
	l.Set("2", Int32(2))
	l.Set("3", Int32(4))
	l.Set("4", Int32(5))
	l.Set("5", Int32(2))
	l.Set("6", Int32(3))

	/*
	 获取用户排名前后2个位置的用户
	*/
	rank := l.GetRank("3")

	left := rank - 2
	right := rank + 2
	if left < 1 {
		left = 1
	}
	if right > l.Len() {
		right = l.Len()
	}
	t.Log(left, rank, right)
	l.Range(left, right, func(rank int, key Key, value interface{}) bool {
		t.Log(rank, "--", key, value)
		return true
	})
}
