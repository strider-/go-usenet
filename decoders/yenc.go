package decoders

import (
	"bufio"
	"bytes"
	"errors"
	"hash/crc32"
	"io"
	"strconv"
	"strings"
)

type Decoder struct {
	r io.Reader
}

type Part struct {
	Crc32, PartCrc32, ExpectedCrc32 uint32
	LineSize                        uint32
	TotalSize, PartSize             uint32
	Part, Total                     uint32
	Begin, End                      uint32
	Name                            string
	Body                            []byte
}

func (d *Decoder) Decode() (*Part, error) {
	esc := new(bool)
	part := &Part{Body: make([]byte, 0)}
	reader := bufio.NewReader(d.r)

	if err := part.readHeader(reader); err != nil {
		return nil, err
	}

	for {
		line, err := readTrimmedLine(reader)
		if err != nil {
			return nil, err
		}

		if len(line) > 5 && string(line[:5]) == "=yend" {
			part.parseYencTrailer(string(line)[6:])
			break
		} else {
			part.Body = append(part.Body, d.decodeLine(line, esc)...)
		}
	}

	part.checksum()

	return part, nil
}

func readTrimmedLine(r *bufio.Reader) ([]byte, error) {
	if line, err := r.ReadBytes('\n'); err != nil {
		return nil, err
	} else {
		return bytes.TrimRight(line, "\r\n"), nil
	}
}

func (d *Decoder) decodeLine(line []byte, esc *bool) []byte {
	c := 0
	for i := 0; i < len(line); i, c = i+1, c+1 {
		if *esc {
			line[c] = (line[i] - 0x40) - 0x2A
			*esc = false
		} else if line[i] == '=' {
			*esc = true
			c--
			continue
		} else {
			line[c] = line[i] - 0x2A
		}
	}
	return line[:c]
}

func (p *Part) ValidPartCrc() bool {
	return p.Crc32 == p.PartCrc32
}

func (p *Part) readHeader(reader *bufio.Reader) error {
	headerline, err := readTrimmedLine(reader)
	if err != nil {
		return err
	}

	if len(headerline) > 7 && string(headerline[:7]) == "=ybegin" {
		p.parseYencHeader(string(headerline)[8:])
	} else {
		return errors.New("missing yenc header")
	}

	if p.Part > 0 {
		partline, err := readTrimmedLine(reader)
		if err != nil {
			return err
		}
		if len(partline) > 6 && string(partline[:6]) == "=ypart" {
			p.parseYencPart(string(partline)[7:])
		} else {
			return errors.New("missing expected yenc part header")
		}
	}
	return nil
}

func (p *Part) parseYencHeader(line string) {
	pairs := p.getKeyValuePairs(line)
	for key, value := range pairs {
		switch key {
		case "part":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.Part = uint32(i)
			}
		case "total":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.Total = uint32(i)
			}
		case "line":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.LineSize = uint32(i)
			}
		case "size":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.TotalSize = uint32(i)
			}
		case "name":
			p.Name = value
		}
	}
}

func (p *Part) parseYencPart(line string) {
	pairs := p.getKeyValuePairs(line)
	for key, value := range pairs {
		switch key {
		case "begin":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.Begin = uint32(i) - 1
			}
		case "end":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.End = uint32(i) - 1
			}
		}
	}
}

func (p *Part) parseYencTrailer(line string) {
	pairs := p.getKeyValuePairs(line)
	for key, value := range pairs {
		switch key {
		case "crc32":
			if i, e := strconv.ParseUint(value, 16, 32); e == nil {
				p.ExpectedCrc32 = uint32(i)
			}
		case "pcrc32":
			if i, e := strconv.ParseUint(value, 16, 32); e == nil {
				p.PartCrc32 = uint32(i)
			}
		case "size":
			if i, e := strconv.ParseUint(value, 10, 32); e == nil {
				p.PartSize = uint32(i)
			}
		}
	}
}

func (p *Part) checksum() {
	hash := crc32.NewIEEE()
	hash.Write(p.Body)
	p.Crc32 = hash.Sum32()
}

func (p *Part) getKeyValuePairs(line string) (result map[string]string) {
	pairs := strings.Split(line, " ")
	result = make(map[string]string)
	for _, kvp := range pairs {
		split := strings.Split(kvp, "=")
		result[split[0]] = split[1]
	}
	return
}

func NewYencDecoder(content []byte) *Decoder {
	return &Decoder{r: bytes.NewReader(content)}
}

func NewYencStreamingDecoder(reader io.Reader) *Decoder {
	return &Decoder{r: reader}
}
