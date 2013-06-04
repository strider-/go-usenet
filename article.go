package nntp

type Article struct {
	Headers map[string]string
	Body    []byte
}
