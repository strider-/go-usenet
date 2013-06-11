package par2

import (
	"bytes"
	"encoding/binary"
	"strings"
)

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
	buff := bytes.NewBuffer(body)
	f.FileID = buff.Next(16)
	f.MD5 = buff.Next(16)
	f.MiniMD5 = buff.Next(16)
	binary.Read(buff, binary.LittleEndian, &f.FileLength)
	f.Filename = strings.TrimRight(string(buff.Next(buff.Len())), "\000")
}
