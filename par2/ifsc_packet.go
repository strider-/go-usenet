package par2

type IFSCPacket struct {
	*Header
	Pairs []ChecksumPair
}

type ChecksumPair struct {
	MD5   []byte
	CRC32 []byte
}

func (i *IFSCPacket) PacketHeader() *Header {
	return i.Header
}

func (i *IFSCPacket) readBody(body []byte) {

}
