package triblab

import (
	//	"fmt"
	"trib"
)

type binClient struct {
	backs    []string
}

func (self *binClient) Bin(name string) trib.Storage {
	h := binHash{backs:self.backs, indexMap: make(map[string]int)}
	h.Init()
	result := &attClient{bin: name, binManager: &h, client: make([]trib.Storage, 3)}
	result.RefreshBin()
	return result
}
