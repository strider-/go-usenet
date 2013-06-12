package par2

import (
	"bytes"
	"encoding/binary"
)

type RecoverySlicePacket struct {
	*Header
	Exponent     uint32
	RecoveryData []byte
}

func (r *RecoverySlicePacket) packetHeader() *Header {
	return r.Header
}

func (r *RecoverySlicePacket) readBody(body []byte) {
	buff := bytes.NewBuffer(body)
	binary.Read(buff, binary.LittleEndian, &r.Exponent)

	r.RecoveryData = buff.Next(int(r.Header.Length) - 4)
}

func (r *RecoverySlicePacket) AvailableBlocks(blocksize uint64) uint64 {
	return uint64(len(r.RecoveryData)) / blocksize
}
