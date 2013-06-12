package par2

import (
	"encoding/binary"
	"io"
)

const validSequence string = "PAR2\000PKT"

type Header struct {
	Sequence      []byte
	Length        uint64
	PacketMD5     []byte
	RecoverySetId []byte
	Type          []byte
	Damaged       bool
}

func (h *Header) fill(r io.Reader) error {
	h.Sequence = make([]byte, 8)
	h.PacketMD5 = make([]byte, 16)
	h.RecoverySetId = make([]byte, 16)
	h.Type = make([]byte, 16)

	_, err := r.Read(h.Sequence)
	if err != nil {
		return err
	}

	binary.Read(r, binary.LittleEndian, &h.Length)
	r.Read(h.PacketMD5)
	r.Read(h.RecoverySetId)
	r.Read(h.Type)
	return nil
}

func (h *Header) ValidSequence() bool {
	return string(h.Sequence) == validSequence
}
