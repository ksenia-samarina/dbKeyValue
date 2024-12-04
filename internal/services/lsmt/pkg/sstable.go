package pkg

import (
	"github.com/bits-and-blooms/bloom"
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/ksenia-samarina/dbKeyValue/internal/services/lsmt/internal"
	"io"
	"os"
)

type SSTable struct {
	index          *redblacktree.Tree
	readBufferSize int
	bloomFilter    *bloom.BloomFilter
	filepath       string
}

func (s *SSTable) Get(key string) (value []byte, ok bool) {
	// first check bloom filter: if key is not present in current sstable
	ok = s.bloomFilter.Test([]byte(key))
	if !ok {
		return nil, false
	}

	node, ok := s.index.Floor(key)
	if !ok {
		return nil, false
	}
	offset := node.Value

	file, err := os.OpenFile(s.filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	_, err = file.Seek(offset.(int64), io.SeekStart)
	if err != nil {
		return nil, false
	}
	reader := internal.NewReader(file)
	for entry, err := reader.ReadEntry(); err == nil; entry, err = reader.ReadEntry() {
		if entry.Key == key {
			return entry.BValue, true
		}
	}
	return nil, false
}

func (s *SSTable) Delete(key string) (ok bool) {
	ok = s.bloomFilter.Test([]byte(key))
	if !ok {
		return false
	}
	// delete key from bloom filter if presented and rebuild it
	ok = s.buildBloomFilter(key)
	return ok
}

func (s *SSTable) buildBloomFilter(excludedKey string) bool {
	file, err := os.OpenFile(s.filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false
	}
	defer file.Close()

	s.bloomFilter.ClearAll()

	deleted := false
	reader := internal.NewReader(file)
	for entry, err := reader.ReadEntry(); err == nil; entry, err = reader.ReadEntry() {
		if entry.Key == excludedKey {
			// key was presented
			deleted = true
			continue
		}
		s.bloomFilter.Add([]byte(entry.Key))
	}
	return deleted
}

func (s *SSTable) buildSparseIndex() {
	file, err := os.OpenFile(s.filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	offset := 0
	prevOffset := 0

	s.bloomFilter.ClearAll()

	reader := internal.NewReader(file)
	for entry, err := reader.ReadEntry(); err == nil; entry, err = reader.ReadEntry() {
		// add each entry key to bloom filter
		s.bloomFilter.Add([]byte(entry.Key))
		if s.index.Empty() || offset-prevOffset > s.readBufferSize {
			s.index.Put(entry.Key, int64(offset))
			prevOffset = offset
		}
		offset += entry.Size()
	}
}

func NewSSTable(filepath string, readBufferSize int, n uint, fp float64) *SSTable {
	s := SSTable{
		index:          redblacktree.NewWithStringComparator(),
		filepath:       filepath,
		readBufferSize: readBufferSize,
		bloomFilter:    bloom.NewWithEstimates(n, fp),
	}
	s.buildSparseIndex()
	return &s
}
