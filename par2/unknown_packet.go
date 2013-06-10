package par2

type UnknownPacket struct {
	*Header
	Body []byte
}

func (u *UnknownPacket) PacketHeader() *Header {
	return u.Header
}

func (u *UnknownPacket) readBody(body []byte) {
	u.Body = body
}
