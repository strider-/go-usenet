package par2

type CreatorPacket struct {
	*Header
	Creator string
}

func (c *CreatorPacket) PacketHeader() *Header {
	return c.Header
}

func (c *CreatorPacket) readBody(body []byte) {

}
