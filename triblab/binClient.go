package triblab

import (
    "trib"
    "hash/fnv"
)

type binClient struct {
    backs []string
    indexMap map[string]trib.Storage
}

func hash(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}

func (self *binClient) Bin(name string) trib.Storage {
    if self.backs == nil {
        return nil
    }
    if st, ok := self.indexMap[name]; ok {
        return st
    } else {
        addr := self.backs[hash(name) % uint32(len(self.backs))]
        self.indexMap[name] = &attClient{bin:name, client:NewClient(addr)}
        return self.indexMap[name]
    }
}