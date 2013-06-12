package par2

import "strings"

type CreatorPacket struct {
	*Header
	Creator string
}

func (c *CreatorPacket) packetHeader() *Header {
	return c.Header
}

func (c *CreatorPacket) readBody(body []byte) {
	c.Creator = strings.TrimRight(string(body), "\000")
}
