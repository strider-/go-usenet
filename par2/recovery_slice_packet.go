package par2

import (
	"bytes"
	"encoding/binary"
)

type RecoverySlicePacket struct {
	*Header
	Exponent     uint32
	RecoveryData [][]byte
}

func (r *RecoverySlicePacket) PacketHeader() *Header {
	return r.Header
}

func (r *RecoverySlicePacket) readBody(body []byte) {
	r.RecoveryData = make([][]byte, 0)
	buff := bytes.NewBuffer(body)
	binary.Read(buff, binary.LittleEndian, &r.Exponent)

	data_len := buff.Len() / 4
	for i := 0; i < data_len; i++ {
		r.RecoveryData = append(r.RecoveryData, buff.Next(4))
	}
}
