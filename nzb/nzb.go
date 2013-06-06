package nzb

import (
	"encoding/xml"
	"io/ioutil"
	"nntp/decoders"
	"time"
)

type Nzb struct {
	name  xml.Name `xml:"nzb"`
	Head  []Meta   `xml:"head>meta"`
	Files []File   `xml:"file"`
}

func (n *Nzb) GenerateQueue(status string) []*QueueItem {
	result := make([]*QueueItem, 0)
	if len(status) == 0 {
		status = "Queued"
	}

	for _, file := range n.Files {
		for _, seg := range file.Segments {
			result = append(result, &QueueItem{file.Subject, seg.Number, uint32(len(file.Segments)), seg.MessageId, status, 1, nil})
		}
	}
	return result
}

type Meta struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",innerxml"`
}

type File struct {
	Poster   string    `xml:"poster,attr"`
	Date     uint64    `xml:"date,attr"`
	Subject  string    `xml:"subject,attr"`
	Groups   []string  `xml:"groups>group"`
	Segments []Segment `xml:"segments>segment"`
}

func (f *File) ParsedDate() time.Time {
	return time.Unix(int64(f.Date), 0)
}

type Segment struct {
	Bytes     uint32 `xml:"bytes,attr"`
	Number    uint32 `xml:"number,attr"`
	MessageId string `xml:",innerxml"`
}

type QueueItem struct {
	FileSubject   string
	SegmentNumber uint32
	TotalSegments uint32
	MessageId     string
	Status        string
	Attempts      int
	Part          *decoders.Part
}

func ReadNzb(filename string) (*Nzb, error) {
	nzb := &Nzb{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err = xml.Unmarshal(file, nzb); err != nil {
		return nil, err
	}

	return nzb, nil
}
