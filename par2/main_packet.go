package par2

type MainPacket struct {
	*Header
	SliceSize             uint64
	RecoverySetCount      uint32
	RecoverySetFileIDs    [][]byte
	NonRecoverySetFileIDs [][]byte
}

func (m *MainPacket) PacketHeader() *Header {
	return m.Header
}

func (m *MainPacket) readBody(body []byte) {

}
