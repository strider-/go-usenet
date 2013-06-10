package par2

type FileDescPacket struct {
	*Header
	FileID     []byte
	MD5        []byte
	MiniMD5    []byte
	FileLength uint64
	Filename   string
}

func (f *FileDescPacket) PacketHeader() *Header {
	return f.Header
}

func (f *FileDescPacket) readBody(body []byte) {

}
