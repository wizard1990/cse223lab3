package triblab

import (
    "hash/fnv"
    "sort"
    "trib"
)

type binHash struct {
    backs    []string
    indexMap map[string]int
}

type ByHashValue []string

func (a ByHashValue) Len() int {
    return len(a)
}

func (a ByHashValue) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a ByHashValue) Less(i, j int) bool {
    return hash(a[i]) < hash(a[j])
}

func (self *binHash) Init() {
    sort.Sort(ByHashValue(self.backs))
}

func hash(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}

func (self *binHash) GetBinCopies(name string) []trib.Storage {
    var startIndex int
    size := len(self.backs)
    if index, ok := self.indexMap[name]; ok {
        startIndex = index
    } else {
        hashValue := hash(name)
        for startIndex = 0; startIndex < size; startIndex++ {
            if hash(self.backs[startIndex]) > hashValue {
                startIndex--
                break
            }
        }
        if startIndex < 0 || startIndex >= size {
            startIndex = size - 1
        }
        self.indexMap[name] = startIndex
    }

    count := 0
    result := make([]trib.Storage, 3)
    for count < 3 {
        if result[count] = self.checkAddr(self.backs[startIndex]); result[count] != nil {
            count++
        }
        startIndex++
        if startIndex >= size {
            startIndex = 0
        }
    }
    return result
}

func (self *binHash) checkAddr(addr string) trib.Storage {
    client := NewClient(addr)
    list := trib.List{[]string{}}
    err := client.Keys(&trib.Pattern{"", "Completed"}, &list)
    if err == nil {
        return client
    } else {
        return nil
    }
}