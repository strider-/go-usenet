package nntp

import (
	"time"
)

type Overview struct {
	ArticleId  uint64
	Subject    string
	Author     string
	Date       time.Time
	MessageId  string
	References string
	Bytes      uint64
	Lines      uint64
}
