package triblab

import (
    "trib"
    "hash/fnv"
    "sort"
)

type binClient struct {
    backs []string
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

func (self *binClient) Init() {
    sort.Sort(ByHashValue(self.backs))
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
        // self.indexMap[name] = &attClient{bin:name, client:NewClient(addr)}
        // return self.indexMap[name]
    }

    count := 0
    result := &attClient{bin:name, client:make([]trib.Storage, 3)}
    for ;count < 3; {
        if result.client[count] = self.checkAddr(self.backs[startIndex]); result.client[count] != nil {
             count++
        }
        startIndex++
        if startIndex >= size {
            startIndex = 0
        }
    }
    return result
}

func (self *binClient) checkAddr(addr string) trib.Storage {
    client := NewClient(addr)
    list := trib.List{[]string{}}
    err := client.Keys(&trib.Pattern{"", "Completed"}, &list)
    if err == nil {
        return client
    } else {
        return nil
    }
}