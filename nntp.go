package nntp

import (
	"crypto/tls"
	"net/textproto"
	"strings"
)

func Dial(addr string, ssl bool) (*Conn, error) {
	conn := new(Conn)

	if ssl {
		if baseConn, err := tls.Dial("tcp", addr, nil); err != nil {
			return nil, err
		} else {
			conn.baseConn = textproto.NewConn(baseConn)
		}
	} else {
		if c, err := textproto.Dial("tcp", addr); err != nil {
			return nil, err
		} else {
			conn.baseConn = c
		}
	}

	if _, _, err := conn.baseConn.ReadCodeLine(200); err != nil {
		return nil, err
	}

	return conn, nil
}

func stripBrackets(mid string) string {
	return strings.NewReplacer("<", "", ">", "").Replace(mid)
}
