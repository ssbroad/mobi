package records

import (
	"encoding/binary"
	"io"

	"github.com/ssbroad/mobi/pdb"
	t "github.com/ssbroad/mobi/types"
)

type EXTHSection struct {
	entries []EXTHEntry
}

func NewEXTHSection() EXTHSection {
	return EXTHSection{
		entries: []EXTHEntry{},
	}
}

func (e *EXTHSection) AddString(tp t.EXTHEntryType, ss ...string) {
	for _, s := range ss {
		if len(s) > 0 {
			entry := NewEXTHEntry(tp, []byte(s))
			e.entries = append(e.entries, entry)
		}
	}
}

func (e *EXTHSection) AddInt(tp t.EXTHEntryType, is ...int) {
	for _, i := range is {
		data := make([]byte, 4)
		pdb.Endian.PutUint32(data, uint32(i))
		entry := NewEXTHEntry(tp, data)
		e.entries = append(e.entries, entry)
	}
}

func (e EXTHSection) Write(w io.Writer) error {
	lenNoPadding := e.LengthWithoutPadding()

	// Write fixed start of header
	h := t.NewEXTHHeader(uint32(lenNoPadding), uint32(len(e.entries)))
	err := binary.Write(w, pdb.Endian, h)
	if err != nil {
		return err
	}

	// Write entries
	for _, entry := range e.entries {
		err := entry.Write(w)
		if err != nil {
			return err
		}
	}

	// Write padding
	pad := make([]byte, invMod(lenNoPadding, 4))
	_, err = w.Write(pad)
	return err
}

func (e EXTHSection) LengthWithoutPadding() int {
	length := t.EXTHHeaderLength
	for _, entry := range e.entries {
		length += entry.Length()
	}

	return length
}

func (e EXTHSection) Length() int {
	length := e.LengthWithoutPadding()
	return length + invMod(length, 4)
}

type EXTHEntry struct {
	EntryType t.EXTHEntryType
	Data      []byte
}

func NewEXTHEntry(tp t.EXTHEntryType, data []byte) EXTHEntry {
	return EXTHEntry{
		EntryType: tp,
		Data:      data,
	}
}

func (e EXTHEntry) Write(w io.Writer) error {
	h := t.NewEXTHEntryHeader(e.EntryType, uint32(e.Length()))
	err := binary.Write(w, pdb.Endian, h)
	if err != nil {
		return err
	}

	_, err = w.Write(e.Data)
	return err
}

func (e EXTHEntry) Length() int {
	return len(e.Data) + t.EXTHEntryHeaderLength
}
