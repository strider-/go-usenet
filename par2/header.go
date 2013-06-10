package par2

var validSequence []byte = []byte{'P', 'A', 'R', '2', 0, 'P', 'K', 'T'}

type Header struct {
	Sequence      []byte
	Length        uint64
	PacketMD5     []byte
	RecoverySetId []byte
	Type          []byte
	Damaged       bool
}

func (h *Header) init() {
	h.Sequence = make([]byte, 8)
	h.PacketMD5 = make([]byte, 16)
	h.RecoverySetId = make([]byte, 16)
	h.Type = make([]byte, 16)
}

func (h *Header) ValidSequence() bool {
	return string(h.Sequence) == string(validSequence)
}
