package par2

type UnknownPacket struct {
	*Header
	Body []byte
}

func (u *UnknownPacket) packetHeader() *Header {
	return u.Header
}

func (u *UnknownPacket) readBody(body []byte) {
	u.Body = body
}
