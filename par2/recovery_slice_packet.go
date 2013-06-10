package par2

type RecoverySlicePacket struct {
	*Header
	Exponent     uint32
	RecoveryData []byte
}

func (r *RecoverySlicePacket) PacketHeader() *Header {
	return r.Header
}

func (r *RecoverySlicePacket) readBody(body []byte) {

}
