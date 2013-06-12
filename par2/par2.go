package par2

import (
	"bytes"
	"crypto/md5"
	"fmt"
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

type ParInfo struct {
	Main         *MainPacket
	Creator      *CreatorPacket
	Files        []*File
	RecoveryData []*RecoverySlicePacket
	ParFiles     []string
	BlockCount   uint32
	TotalSize    uint64
}

func Stat(file string) (*ParInfo, error) {
	par_files, err := allParFiles(file)
	if err != nil {
		return nil, err
	}

	stat := &ParInfo{nil, nil, make([]*File, 0), make([]*RecoverySlicePacket, 0), par_files, 0, 0}
	packets, err := packets(stat.ParFiles)
	if err != nil {
		return nil, err
	}

	table := make(map[string]*File)
	for _, p := range packets {
		switch p.(type) {
		case *MainPacket:
			stat.Main = p.(*MainPacket)
		case *CreatorPacket:
			stat.Creator = p.(*CreatorPacket)
		case *RecoverySlicePacket:
			stat.RecoveryData = append(stat.RecoveryData, p.(*RecoverySlicePacket))
		case *FileDescPacket:
			tmp := p.(*FileDescPacket)
			id := fmt.Sprintf("%x", tmp.FileID)
			stat.TotalSize += tmp.FileLength
			if val, exists := table[id]; exists {
				val.FileDescPacket = tmp
				stat.Files = append(stat.Files, val)
			} else {
				table[id] = &File{tmp, nil}
			}
		case *IFSCPacket:
			tmp := p.(*IFSCPacket)
			id := fmt.Sprintf("%x", tmp.FileID)
			stat.BlockCount += uint32(len(tmp.Pairs))
			if val, exists := table[id]; exists {
				val.IFSCPacket = tmp
				stat.Files = append(stat.Files, val)
			} else {
				table[id] = &File{nil, tmp}
			}
		}
	}

	return stat, nil
}

func allParFiles(file string) ([]string, error) {
	dir, fname := filepath.Split(file)
	ext := filepath.Ext(fname)
	return filepath.Glob(dir + fname[:len(fname)-len(ext)] + ".*par2")
}

func packets(files []string) ([]Packet, error) {
	packets := make([]Packet, 0)
	for _, par := range files {
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
