package records

import (
	"encoding/binary"
	"io"

	"github.com/ssbroad/mobi/pdb"
	t "github.com/ssbroad/mobi/types"
)

type IndexRecord struct {
	Type          uint32
	HeaderType    uint32
	TAGXTable     t.TAGXTagTable
	IDXTEntries   [][]byte
	SubEntryCount uint32
	CNCXCount     uint32
}

func (r IndexRecord) Write(w io.Writer) error {
	// Headers
	inh := t.NewINDXHeader(0, 0)
	th := t.NewTAGXHeader()
	idh := t.NewIDXTHeader()
	offset := t.INDXHeaderLength

	// INDX variables
	inh.IndexRecordCount = uint32(len(r.IDXTEntries))
	inh.IndexType = r.Type
	inh.HeaderType = r.HeaderType
	inh.CNCXCount = r.CNCXCount

	// TAGX variables
	if len(r.TAGXTable) > 0 {
		inh.TAGXOffset = t.INDXHeaderLength
		th.HeaderLength += uint32(len(r.TAGXTable) * t.TAGXTagLength)
		offset += int(th.HeaderLength)
	} else {
		inh.TAGXOffset = 0
	}

	// IDXT variables
	idxtLength := 0
	idxtOffsets := make([]uint16, 0)
	for _, entry := range r.IDXTEntries {
		length := len(entry)
		idxtOffsets = append(idxtOffsets, uint16(offset))
		offset += length
		idxtLength += length
	}
	inh.IDXTStart = uint32(offset + idxtLength%4)
	inh.IndexEntryCount = r.SubEntryCount

	// Write INDX header
	err := binary.Write(w, pdb.Endian, inh)
	if err != nil {
		return err
	}

	if len(r.TAGXTable) > 0 {
		// Write TAGX section
		err := writeSequential(w, pdb.Endian, th, r.TAGXTable)
		if err != nil {
			return err
		}
	}

	// Write entries
	for _, entry := range r.IDXTEntries {
		err := binary.Write(w, pdb.Endian, entry)
		if err != nil {
			return err
		}
	}

	// Write IDXT and padding
	idxtPad := make([]byte, idxtLength%4)
	postPad := make([]byte, r.LengthNoPadding()%4)
	return writeSequential(w, pdb.Endian, idxtPad, idh, idxtOffsets, postPad)
}

func (r IndexRecord) Length() int {
	length := r.LengthNoPadding()
	return length + length%4
}

func (r IndexRecord) LengthNoPadding() int {
	length := t.INDXHeaderLength + t.IDXTHeaderLength + len(r.IDXTEntries)*2
	entriesLength := 0
	for _, entry := range r.IDXTEntries {
		length += len(entry)
		entriesLength += len(entry)
	}
	length += entriesLength % 4
	if len(r.TAGXTable) > 0 {
		length += t.TAGXHeaderLength + len(r.TAGXTable)*t.TAGXTagLength
	}

	return length
}
