# sortedset

- 自定义排序规则
- 遍历有序集合中指定区间分数的成员
- 通过索引区间返回有序集合指定区间内的成员
- 获取排名前 N 的用户、翻页
- 获取用户排名附近的用户


```
key string , value 为实现排序接口的任意类型

type Key string
type Interface interface {
	Less(other interface{}) bool
}
```

使用 map 和 skiplist 实现

skiplist 跳表，它允许快速查询，插入和删除一个有序连续元素的数据链表。
跳跃列表的平均查找和插入时间复杂度都是O(logn)。

## 用法

1. 获取用户附近排名. `radius = 10`
    1) 通过 `GetRank` 获取当前用户的排名 rank。
    2) 计算出区间 `[rank-radius: rank+radius]`,通过 `Range` 遍历。
    
2. 获取排名前 N 的用户. `zs.Range(1, N)`, 亦可实现翻页。

3. 获取分数区间 `[s1:s2]` 的用户数量、用户。`range by score`.
    1) 通过 `Search` 函数分别获取 s1,s2 的排名 r1,r2 ,仅返回将要插入的位置。
    2) 计算用户数量： r1-r2 取绝对值。
    3) 遍历区间用户： 通过 `Range` 遍历。
   
4. 指定用户分数加上增量。`increase by`.
    1) 获取用户对象后移除
    2) 用户对象加上分数后添加

5. 获取大于、小于 `score` 的用户数量。`zs.Search` 实现。


## Usage

```go
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

	t.Log(zs.Len())                // 6
	// get rank by key
	t.Log(zs.GetRank("hello"))     // 5
	// get value by key 
	t.Log(zs.GetValue("hello"))    // 2.2 true  
	// get key,value by rank
	t.Log(zs.Select(zs.Len()))     // , 1.1

	// update
	fmt.Println()
	zs.Set("hello", Score(6.6))
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
	t.Log(zs.Search(zs.Len(), func(i Interface) bool {  // 3
		return i.(Score) > score
	}))

	// would be inserted
	t.Log(zs.WouldBeInserted(score))   // 5

}
```
```go
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

//    output
//    sortedset_test.go:92: 1  --  u2 &{u2 2 40}
//    sortedset_test.go:92: 2  --  u4 &{u4 3 30}
//    sortedset_test.go:92: 3  --  u3 &{u3 3 30}
//    sortedset_test.go:92: 4  --  u1 &{u1 2 30}
```
