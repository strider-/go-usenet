package par2

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"os"
)

const (
	headerLength            uint64 = 0x40
	typeMainPacket          string = "PAR 2.0\000Main\000\000\000\000"
	typeFileDescPacket      string = "PAR 2.0\000FileDesc"
	typeIFSCPacket          string = "PAR 2.0\000IFSC\000\000\000\000"
	typeRecoverySlicePacket string = "PAR 2.0\000RecvSlic"
	typeCreatorPacket       string = "PAR 2.0\000Creator\000"
)

type Packet interface {
	readBody([]byte)
	PacketHeader() *Header
}

func Packets(file string) ([]Packet, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	packets := make([]Packet, 0)

	for {
		h := new(Header)
		h.init()

		buf := make([]byte, 8)

		_, err := f.Read(h.Sequence)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		f.Read(buf)
		binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, &h.Length)
		f.Read(h.PacketMD5)
		f.Read(h.RecoverySetId)
		f.Read(h.Type)

		buf = make([]byte, h.Length-headerLength)
		f.Read(buf)

		p := createPacket(h)
		verifyPacket(h, buf)
		p.readBody(buf)

		packets = append(packets, p)
	}

	return packets, nil
}

func createPacket(h *Header) Packet {
	switch string(h.Type) {
	case typeMainPacket:
		return &MainPacket{h, 0, 0, nil, nil}
	case typeFileDescPacket:
		return &FileDescPacket{h, nil, nil, nil, 0, ""}
	case typeIFSCPacket:
		return &IFSCPacket{h, nil, nil}
	case typeRecoverySlicePacket:
		return &RecoverySlicePacket{h, 0, nil}
	case typeCreatorPacket:
		return &CreatorPacket{h, ""}
	}

	return &UnknownPacket{h, nil}
}

func verifyPacket(h *Header, body []byte) {
	hash := md5.New()
	hash.Write(h.RecoverySetId)
	hash.Write(h.Type)
	hash.Write(body)

	h.Damaged = (len(body)%4) != 0 || !bytes.Equal(hash.Sum(nil), h.PacketMD5)
}
