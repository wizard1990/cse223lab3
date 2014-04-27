package ref

import (
	"trib"
)

type seqTrib struct {
	seq uint64
	*trib.Trib
}

type bySeq []*seqTrib

func (s bySeq) Len() int           { return len(s) }
func (s bySeq) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySeq) Less(i, j int) bool { return s[i].seq < s[j].seq }
