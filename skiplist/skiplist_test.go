package skiplist

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type User struct {
	name  string
	score int
}

func (this *User) Less(other interface{}) bool {
	return this.score >= other.(*User).score
}

func TestNew(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	l := New()

	e, r1 := l.Insert(&User{name: "1", score: 1})
	l.Insert(&User{name: "4", score: 4})
	l.Insert(&User{name: "6", score: 6})
	l.Insert(&User{name: "3", score: 3})
	l.Insert(&User{name: "5", score: 5})
	l.Insert(&User{name: "2", score: 2})
	l.Insert(&User{name: "33", score: 3})
	e2, r2 := l.Insert(&User{name: "11", score: 1})

	t.Log(l.Len(), r1, r2)
	for e, i := l.Front(), 1; e != nil; e, i = e.Next(), i+1 {
		u := e.Value().(*User)
		t.Logf("%d, user(%s,%d)", i, u.name, u.score)
	}

	//for e, i := l.Back(), 1; e != nil; e, i = e.Prev(), i+1 {
	//	u := e.Value().(*User)
	//	t.Logf("%d, user(%s,%d)", i, u.name, u.score)
	//}
	//printSl(l)
	//t.Log(e.Rank(), e2.Rank(), l.GetRank(e2))
	t.Log()

	l.Remove(e)
	t.Log(l.Len(), e.Rank(), e2.Rank(), l.GetRank(e2))
	e, r1 = l.Insert(e.Value().(*User))
	t.Log(l.Len(), e.Rank(), e2.Rank(), l.GetRank(e2))

	for e, i := l.Front(), 1; e != nil; e, i = e.Next(), i+1 {
		u := e.Value().(*User)
		t.Logf("%d, user(%s,%d)", i, u.name, u.score)
	}

	t.Log(l.Select(1).Value(), l.Select(5).Value())

	// Search
	// find the smallest rank in descending order.
	u := &User{score: 3}
	f := func(i Interface) bool {
		return i.(*User).score > u.score
	}
	t.Log(l.Search(l.Len(), f))

	// WouldBeInserted
	t.Log(l.WouldBeInserted(&User{score: 5}))
}

func printSl(sl *SkipList) {
	for i := sl.level - 1; i >= 0; i-- {
		str := []string{}
		for e := sl.head; e != sl.tail; e = e.links[i].next {
			str = append(str, fmt.Sprintf("-%d->%v", e.links[i].skip, e.links[i].next.value))
		}
		fmt.Println(strings.Join(str, " "))
	}

	str := ""
	for e, i := sl.Front(), 1; e != nil; e, i = e.Next(), i+1 {
		str += fmt.Sprintf("%v  ", e.value)
	}
	fmt.Println(str)

}
