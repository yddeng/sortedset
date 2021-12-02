package skiplist

import (
	"fmt"
	"testing"
)

type User struct {
	score float64
	id    string
}

func (u *User) Less(other interface{}) bool {
	if u.score > other.(*User).score {
		return true
	}
	if u.score == other.(*User).score && len(u.id) > len(other.(*User).id) {
		return true
	}
	return false
}

func TestNew(t *testing.T) {
	us := make([]*User, 7)
	us[0] = &User{6.6, "hi"}
	us[1] = &User{4.4, "hello"}
	us[2] = &User{2.2, "world"}
	us[3] = &User{3.3, "go"}
	us[4] = &User{1.1, "skip"}
	us[5] = &User{2.2, "list"}
	us[6] = &User{3.3, "lang"}

	// insert
	sl := New()
	t.Log(sl.Len(), sl.Front(), sl.Back())

	for i := 0; i < len(us); i++ {
		sl.Insert(us[i])
	}
	t.Log(sl.Len(), sl.Front().Value, sl.Back().Value)

	// traverse
	for e := sl.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*User).id, "-->", e.Value.(*User).score)
	}
	t.Log(sl.Len())

	// rank
	rank1 := sl.GetRank(&User{2.2, "list"})
	rank2 := sl.GetRank(&User{6.6, "hi"})
	if rank1 != 6 || rank2 != 1 {
		t.Fatal()
	}
	if e := sl.GetElementByRank(2); e.Value.(*User).score != 4.4 || e.Value.(*User).id != "hello" {
		t.Fatal(e)
	}

	nu := &User{score: 0, id: "dfsf"}
	t.Log(sl.Find(nu), sl.GetRank(nu))

	t.Log(sl.GetRank(us[0]), sl.GetElementByRank(sl.Len()).Value)
}
