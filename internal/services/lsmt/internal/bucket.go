package internal

import "github.com/google/btree"

type Bucket struct {
	btree *btree.BTree
	size  uint
}

func (b *Bucket) Put(key string, val []byte) (replaced bool) {
	item := NewNode(key, val)
	i := b.btree.ReplaceOrInsert(item)
	// if item is nil than new node was added
	if i == nil {
		b.size += item.Size()
		return false
	}
	return true
}

func (b *Bucket) Get(key string) (node Node, ok bool) {
	i := b.btree.Get(NewNode(key, []byte{}))
	if i != nil {
		return i.(Node), true
	}
	return NewNode("", []byte{}), false
}

func (b *Bucket) Delete(key string) (ok bool) {
	i := b.btree.Delete(NewNode(key, []byte{}))
	if i != nil {
		return true
	}
	return false
}

func (b *Bucket) Scan() (nodes []Node, ok bool) {
	var items []Node
	next := func(i btree.Item) bool {
		items = append(items, i.(Node))
		return true
	}
	b.btree.Ascend(next)
	if len(items) > 0 {
		return items, true
	}
	return nil, false
}

func (b *Bucket) Size() uint {
	return b.size
}

func NewBucket(n int) *Bucket {
	return &Bucket{
		btree: btree.New(n),
	}
}
