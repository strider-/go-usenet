package nzb

import (
	"encoding/xml"
	"io/ioutil"
	"time"
)

type Nzb struct {
	name  xml.Name `xml:"nzb"`
	Head  []Meta   `xml:"head>meta"`
	Files []File   `xml:"file"`
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

func OpenNzb(filename string) (*Nzb, error) {
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
