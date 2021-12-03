package grank

import "github.com/yddeng/grank/skiplist"

type skiplistSlot struct {
	list   *skiplist.SkipList
	header *skiplist.Element
	tail   *skiplist.Element
}

type SkiplistBuckup struct {
	slotList []*skiplistSlot
}

func newSkiplistSlot() *skiplistSlot {
	list := skiplist.New()
	return &skiplistSlot{
		list:   list,
		header: list.Front(),
		tail:   list.Back(),
	}
}
