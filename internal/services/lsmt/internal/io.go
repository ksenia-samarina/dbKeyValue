package internal

import (
	"encoding/binary"
	"os"
)

type Reader struct {
	file *os.File
}

func (r *Reader) ReadEntry() (entry *Entry, err error) {
	keyLengthByteArray := make([]byte, 8)
	_, err = r.file.Read(keyLengthByteArray)
	if err != nil {
		return nil, err
	}
	keyLength := binary.BigEndian.Uint64(keyLengthByteArray)

	valueLengthByteArray := make([]byte, 8)
	_, err = r.file.Read(valueLengthByteArray)
	if err != nil {
		return nil, err
	}
	valueLength := binary.BigEndian.Uint64(valueLengthByteArray)

	keyArray := make([]byte, keyLength)
	_, err = r.file.Read(keyArray)
	if err != nil {
		return nil, err
	}
	key := string(keyArray)

	value := make([]byte, valueLength)
	_, err = r.file.Read(value)
	if err != nil {
		return nil, err
	}
	return NewEntry(key, value), nil
}

func NewReader(f *os.File) *Reader {
	return &Reader{file: f}
}
