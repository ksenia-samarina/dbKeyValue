package internal

import "encoding/binary"

type Entry struct {
	Key    string
	BValue []byte
}

func (e *Entry) Marshall() []byte {
	bKey := []byte(e.Key)

	keyLength := make([]byte, 8)
	binary.BigEndian.PutUint64(
		keyLength,
		uint64(len(bKey)),
	)

	valueLength := make([]byte, 8)
	binary.BigEndian.PutUint64(
		valueLength,
		uint64(len(e.BValue)),
	)

	var data []byte
	for _, b := range [][]byte{keyLength, valueLength, bKey, e.BValue} {
		data = append(data, b...)
	}
	return data
}

func (e *Entry) Unmarshall(in []byte) *Entry {
	keyLength := binary.BigEndian.Uint64(in[0:8])
	valueLength := binary.BigEndian.Uint64(in[8:16])

	key := string(in[16:(16 + keyLength)])
	bValue := in[(16 + keyLength):(16 + keyLength + valueLength)]

	return NewEntry(key, bValue)
}

func (e *Entry) Size() int {
	return len(e.Marshall())
}

func NewEntry(key string, value []byte) *Entry {
	return &Entry{
		Key:    key,
		BValue: value,
	}
}
