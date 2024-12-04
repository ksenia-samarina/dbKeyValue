package internal

import "github.com/google/btree"

type Node struct {
	Key string
	Val []byte
}

func (i Node) Less(than btree.Item) bool {
	return i.Key < than.(Node).Key
}

func (i Node) Size() uint {
	return uint(len(i.Val) + len(i.Key))
}

func NewNode(key string, val []byte) Node {
	return Node{
		Key: key,
		Val: val,
	}
}
