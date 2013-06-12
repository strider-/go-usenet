package par2

import (
	"bytes"
	"encoding/binary"
)

type MainPacket struct {
	*Header
	BlockSize             uint64
	RecoverySetCount      uint32
	RecoverySetFileIDs    [][]byte
	NonRecoverySetFileIDs [][]byte
}

func (m *MainPacket) packetHeader() *Header {
	return m.Header
}

func (m *MainPacket) readBody(body []byte) {
	m.RecoverySetFileIDs = make([][]byte, 0)
	m.NonRecoverySetFileIDs = make([][]byte, 0)

	buff := bytes.NewBuffer(body)
	binary.Read(buff, binary.LittleEndian, &m.BlockSize)
	binary.Read(buff, binary.LittleEndian, &m.RecoverySetCount)

	for i := 0; i < int(m.RecoverySetCount); i++ {
		m.RecoverySetFileIDs = append(m.RecoverySetFileIDs, buff.Next(16))
	}

	non_rec_count := buff.Len() / 16
	for i := 0; i < non_rec_count; i++ {
		m.NonRecoverySetFileIDs = append(m.NonRecoverySetFileIDs, buff.Next(16))
	}
}
