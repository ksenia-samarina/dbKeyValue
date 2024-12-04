package pkg

import (
	"github.com/ksenia-samarina/dbKeyValue/internal/services/lsmt/internal"
	"io"
)

type MemTable struct {
	bucket *internal.Bucket

	Timestamp int64
}

func (m *MemTable) Put(key string, value []byte) {
	m.bucket.Put(key, value)
}

func (m *MemTable) Get(key string) (value []byte, ok bool) {
	var node internal.Node
	node, ok = m.bucket.Get(key)
	if ok {
		return node.Val, true
	}
	return nil, false
}

func (m *MemTable) Delete(key string) bool {
	ok := m.bucket.Delete(key)
	if ok {
		return true
	}
	return false
}

func (m *MemTable) Write(w io.Writer) (n int, err error) {
	written := 0
	items, _ := m.bucket.Scan()
	for _, item := range items {
		entry := internal.NewEntry(item.Key, item.Val)
		n, err = w.Write(entry.Marshall())
		if err != nil {
			return 0, err
		}
		written += n
	}
	return written, nil
}

func (m *MemTable) Size() uint {
	return m.bucket.Size()
}

func NewMemTable(n int) *MemTable {
	return &MemTable{
		bucket: internal.NewBucket(n),
	}
}
