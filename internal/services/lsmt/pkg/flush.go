package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Flush struct {
	ssTableFilename string
	memTable        *MemTable
}

func (f *Flush) Flush() string {
	log.Printf("[DEBUG] Starting memtable flushing process for file=%s", f.ssTableFilename)

	file, err := os.OpenFile(f.ssTableFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
	}
	defer file.Close()

	_, err = f.memTable.Write(file)
	if err != nil {

	}
	err = file.Sync()
	if err != nil {
	}
	log.Printf("[DEBUG] memtable saved as SSTable to the file=%s", file.Name())
	return file.Name()
}

func NewFlush(ssTablesDir string, memTable *MemTable) *Flush {
	ssTableFilename := filepath.Join(ssTablesDir, fmt.Sprintf("%v.sstable", memTable.Timestamp))
	return &Flush{
		ssTableFilename: ssTableFilename,
		memTable:        memTable,
	}
}
