package nntp

import (
	"fmt"
	"net/textproto"
	"nntp/decoders"
	"strconv"
	"strings"
	"time"
)

type Conn struct {
	baseConn *textproto.Conn
	group    string
}

func (c *Conn) Authenticate(user, pass string) error {
	if _, _, err := c.sendCmd(fmt.Sprintf("AUTHINFO USER %s", user), 381); err != nil {
		return err
	} else {
		if _, _, err := c.sendCmd(fmt.Sprintf("AUTHINFO PASS %s", pass), 281); err != nil {
			return err
		}
	}

	return nil
}

func (c *Conn) Groups() ([]Group, error) {
	if _, _, err := c.sendCmd("LIST", 215); err != nil {
		return nil, err
	}

	grps, err := c.baseConn.ReadDotLines()
	if err != nil {
		return nil, err
	}

	result := make([]Group, len(grps))
	for i := range grps {
		var h, l uint64
		var err error
		split := strings.Split(grps[i], " ")

		if h, err = strconv.ParseUint(split[1], 10, 64); err != nil {
			return nil, err
		}
		if l, err = strconv.ParseUint(split[2], 10, 64); err != nil {
			return nil, err
		}

		result[i] = Group{Name: split[0], High: h, Low: l, CanPost: split[3] == "y"}
	}

	return result, nil
}

func (c *Conn) SetGroup(group string) (result *Group, err error) {
	var l, h uint64
	var msg string

	if _, msg, err = c.sendCmd(fmt.Sprintf("GROUP %s", group), 211); err != nil {
		return
	}

	split := strings.Split(msg, " ")
	if l, err = strconv.ParseUint(split[1], 10, 64); err != nil {
		return
	}
	if h, err = strconv.ParseUint(split[2], 10, 64); err != nil {
		return
	}

	result = &Group{Name: split[3], High: h, Low: l}
	c.group = result.Name

	return
}

func (c *Conn) Group() string {
	return c.group
}

func (c *Conn) Head(mid string) (article *Article, err error) {
	article = &Article{}
	if _, _, err = c.sendCmd(fmt.Sprintf("HEAD %s", mid), 221); err != nil {
		return
	}

	err = c.fillHeaders(article)
	return
}

func (c *Conn) Article(mid string) (article *Article, err error) {
	article = &Article{}
	mid = stripBrackets(mid)

	if _, _, err = c.sendCmd(fmt.Sprintf("ARTICLE <%s>", mid), 220); err != nil {
		return
	}

	if err = c.fillHeaders(article); err != nil {
		return
	}

	article.Body, err = c.baseConn.ReadDotBytes()
	return
}

func (c *Conn) DecodedArticle(mid string) (*decoders.Part, error) {
	mid = stripBrackets(mid)
	if _, _, err := c.sendCmd(fmt.Sprintf("BODY <%s>", mid), 222); err != nil {
		return nil, err
	}

	dec := decoders.NewYencStreamingDecoder(c.baseConn.R)
	part, err := dec.Decode()
	if err != nil {
		return nil, err
	}
	c.baseConn.ReadDotBytes()
	return part, nil
}

func (c *Conn) Exists(mid string) bool {
	_, _, err := c.sendCmd(fmt.Sprintf("STAT %s", mid), 223)
	return err == nil
}

func (c *Conn) Close() error {
	c.sendCmd("QUIT", 205)
	return c.baseConn.Close()
}

func (c *Conn) Date() (time.Time, error) {
	_, rawDate, err := c.sendCmd("DATE", 111)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), err
	}
	return time.Parse("20060102150405", rawDate)
}

func (c *Conn) sendCmd(cmd string, expectCode int) (int, string, error) {
	id, err := c.baseConn.Cmd(cmd)
	if err != nil {
		return 0, "", err
	}
	c.baseConn.StartResponse(id)
	defer c.baseConn.EndResponse(id)
	return c.baseConn.ReadCodeLine(expectCode)
}

func (c *Conn) fillHeaders(article *Article) error {
	for {
		if cur, err := c.baseConn.ReadLine(); err != nil {
			return err
		} else {
			if cur != "" && cur != "." {
				split := strings.SplitN(cur, ":", 2)
				article.Headers[split[0]] = split[1]
			} else {
				break
			}
		}
	}

	return nil
}
