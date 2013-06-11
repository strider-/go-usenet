package par2

import (
	"bytes"
)

type IFSCPacket struct {
	*Header
	FileID []byte
	Pairs  []ChecksumPair
}

type ChecksumPair struct {
	MD5   []byte
	CRC32 []byte
}

func (i *IFSCPacket) PacketHeader() *Header {
	return i.Header
}

func (i *IFSCPacket) readBody(body []byte) {
	i.Pairs = make([]ChecksumPair, 0)
	buff := bytes.NewBuffer(body)
	i.FileID = buff.Next(16)

	pair_count := buff.Len() / 20
	for n := 0; n < pair_count; n++ {
		pair := &ChecksumPair{buff.Next(16), buff.Next(4)}
		i.Pairs = append(i.Pairs, *pair)
	}
}
