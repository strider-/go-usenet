package par2

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
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

func allParFiles(file string) ([]string, error) {
	dir, fname := filepath.Split(file)
	ext := filepath.Ext(fname)
	return filepath.Glob(dir + fname[:len(fname)-len(ext)] + ".*par2")
}

func Packets(file string) ([]Packet, error) {
	pars, err := allParFiles(file)
	if err != nil {
		return nil, err
	}

	packets := make([]Packet, 0)
	for _, par := range pars {
		f, err := os.Open(par)
		if err != nil {
			return nil, err
		}

		defer f.Close()
		stat, _ := f.Stat()
		par_size := stat.Size()

		for {
			h := new(Header)
			if err := h.fill(f); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			if !h.ValidSequence() {
				r, _ := f.Seek(-7, os.SEEK_CUR)
				if (par_size - r) < 8 {
					break
				}
				continue
			}

			buf := make([]byte, h.Length-headerLength)
			f.Read(buf)

			p := createPacket(h)
			verifyPacket(h, buf)
			p.readBody(buf)

			if !h.Damaged && !contains(packets, p) {
				packets = append(packets, p)
			}
		}
	}

	return packets, nil
}

func contains(packets []Packet, packet Packet) bool {
	header := packet.PacketHeader()
	for _, p := range packets {
		h := p.PacketHeader()
		if string(h.PacketMD5) == string(header.PacketMD5) {
			return true
		}
	}
	return false
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
