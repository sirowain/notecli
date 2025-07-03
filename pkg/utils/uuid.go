package utils

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type NoteId uint64

func (n NoteId) String() string {
	return fmt.Sprintf("%03d", n)
}

func Btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Itob(v NoteId) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func Stoi(s string) (NoteId, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return NoteId(0), err
	}

	return NoteId(id), nil
}
