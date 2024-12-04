package pkg

import (
	"context"
	"sync"
	"time"
)

var flushMutex = &sync.Mutex{}
var ssTablesMutex = &sync.Mutex{}

type LSMT struct {
	memTable            *MemTable
	ssTables            []*SSTable
	memTablesFlushQueue []*MemTable

	maxMemTableSize uint
	btreeSize       int
	ssTablesDir     string

	bloomFilterN          uint
	bloomFilterFp         float64
	ssTableReadBufferSize int

	status int
}

func (l *LSMT) Put(ctx context.Context, key string, val []byte) bool {
	l.flushMemTable()
	l.memTable.Put(key, val)
	return true
}

func (l *LSMT) flushMemTable() {
	if l.memTable.Size() > l.maxMemTableSize {
		memTableCopy := l.memTable
		memTableCopy.Timestamp = time.Now().UnixNano()

		l.memTable = NewMemTable(l.btreeSize)

		go func() {
			flushMutex.Lock()
			defer flushMutex.Unlock()

			l.memTablesFlushQueue = append(l.memTablesFlushQueue, memTableCopy)
		}()
	}
}

func (l *LSMT) Get(ctx context.Context, key string) (val []byte, ok bool) {
	val, ok = l.memTable.Get(key)
	if !ok {
		val, ok = l.getFromMemTablesFlushQueue(key)
	}
	if !ok {
		val, ok = l.getFromSSTables(key)
	}
	return val, ok
}

func (l *LSMT) getFromMemTablesFlushQueue(key string) (val []byte, ok bool) {
	for i := len(l.memTablesFlushQueue) - 1; i >= 0; i-- {
		val, ok = l.memTablesFlushQueue[i].Get(key)
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (l *LSMT) getFromSSTables(key string) (val []byte, ok bool) {
	for i := len(l.ssTables) - 1; i >= 0; i-- {
		val, ok = l.ssTables[i].Get(key)
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (l *LSMT) Delete(ctx context.Context, key string) bool {
	ok := l.memTable.Delete(key)
	if !ok {
		ok = l.deleteFromMemTablesFlushQueue(key)
	}
	if !ok {
		ok = l.deleteFromSSTables(key)
	}
	return ok
}

func (l *LSMT) deleteFromMemTablesFlushQueue(key string) bool {
	deleted := false
	for i := len(l.memTablesFlushQueue) - 1; i >= 0; i-- {
		ok := l.memTablesFlushQueue[i].Delete(key)
		if ok {
			deleted = true
		}
	}
	return deleted
}

func (l *LSMT) deleteFromSSTables(key string) bool {
	deleted := false
	for i := len(l.ssTables) - 1; i >= 0; i-- {
		ok := l.ssTables[i].Delete(key)
		if ok {
			deleted = true
		}
	}
	return deleted
}

func (l *LSMT) startFlushing() {
	for l.status == 0 {
		flushMutex.Lock()
		for _, memTable := range l.memTablesFlushQueue {
			f := NewFlush(l.ssTablesDir, memTable)
			filepath := f.Flush()
			ssTablesMutex.Lock()
			ssTable := NewSSTable(filepath, l.ssTableReadBufferSize, l.bloomFilterN, l.bloomFilterFp)
			l.ssTables = append(l.ssTables, ssTable)
			ssTablesMutex.Unlock()
		}
		l.memTablesFlushQueue = []*MemTable{}
		flushMutex.Unlock()
		time.Sleep(time.Millisecond * 100)
	}
}

func NewLSMT(btreeSize int, maxMemTableSize uint, ssTablesDir string, bloomFilterN uint, bloomFilterFp float64) *LSMT {
	l := LSMT{
		memTable:            NewMemTable(btreeSize),
		ssTables:            make([]*SSTable, 0),
		memTablesFlushQueue: make([]*MemTable, 0),
		maxMemTableSize:     maxMemTableSize,
		btreeSize:           btreeSize,
		ssTablesDir:         ssTablesDir,
		bloomFilterN:        bloomFilterN,
		bloomFilterFp:       bloomFilterFp,
	}
	l.status = 0
	go l.startFlushing()
	return &l
}
