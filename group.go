package nntp

type Group struct {
	Name      string
	High, Low uint64
	CanPost   bool
}

func (g Group) Count() uint64 {
	return g.High - g.Low
}
