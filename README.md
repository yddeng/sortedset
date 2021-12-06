# sortedset

自定义排序规则、多条件排序。

```
type Key string
type Interface interface {
	Less(other interface{}) bool
}
```

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
	t.Log(zs.GetByRank(zs.Len()))  // , 1.1

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
	t.Log(zs.Delete("hello"))       // 6.6
	// RevRange
	rank = zs.Len()
	zs.RevRange(1, zs.Len(), func(key Key, value interface{}) {
		t.Log(rank, " -- ", key, value.(Score))
		rank--
	})
}
```
```
    sortedset_test.go:24: 6
    sortedset_test.go:26: 5
    sortedset_test.go:28: 2.2 true
    sortedset_test.go:30: , 1.1

    sortedset_test.go:38: 1  --  hello 6.6
    sortedset_test.go:38: 2  --  world 5.5
    sortedset_test.go:38: 3  --  you 5.5
    sortedset_test.go:38: 4  --  how 3.3
    sortedset_test.go:38: 5  --  are 3.3
    sortedset_test.go:38: 6  --  , 1.1

    sortedset_test.go:44: 6.6
    sortedset_test.go:48: 5  --  , 1.1
    sortedset_test.go:48: 4  --  are 3.3
    sortedset_test.go:48: 3  --  how 3.3
    sortedset_test.go:48: 2  --  you 5.5
    sortedset_test.go:48: 1  --  world 5.5
```


## skiplist

跳表，跳表全称为跳跃列表，它允许快速查询，插入和删除一个有序连续元素的数据链表。
跳跃列表的平均查找和插入时间复杂度都是O(logn)。

快速查询是通过维护一个多层次的链表，且每一层链表中的元素是前一层链表元素的子集。
一开始时，算法在最稀疏的层次进行搜索，直至需要查找的元素在该层两个相邻的元素中间。
这时，算法将跳转到下一个层次，重复刚才的搜索，直到找到需要查找的元素为止。

## 拓展

1. 获取用户附近排名. radius = 10
    1) 通过 `GetRank` 获取当前用户的排名 rank.
    2) 计算出区间 `[rank-radius: rank+radius]`,通过 `Range` 遍历
    
2. 获取排名前 N 的用户. `zs.Range(1, N)`, 可实现翻页.
    
