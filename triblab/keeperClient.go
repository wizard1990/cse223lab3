package triblab

import (
	"trib"
	//"trib/colon"
)

type keeperClient struct {
	bin    string
	client trib.Storage
}

/*
func genPrefix(s string) string {
	return colon.Escape(s) + "::"
}
*/

type KeeperStorage interface {
	KeeperBin(name string) []trib.Storage
}

func (self *keeperClient) Get(key string, value *string) error {
	return self.client.Get(genPrefix(self.bin)+key, value)
}

func (self *keeperClient) Set(kv *trib.KeyValue, succ *bool) error {
	return self.client.Set(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, succ)
}

func (self *keeperClient) Keys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix}
	if e := self.client.Keys(&np, list); e != nil {
		return e
	}
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
	}
	return nil
}

func (self *keeperClient) ListGet(key string, list *trib.List) error {
	return self.client.ListGet(genPrefix(self.bin)+key, list)
}

func (self *keeperClient) ListAppend(kv *trib.KeyValue, succ *bool) error {
	return self.client.ListAppend(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, succ)
}

func (self *keeperClient) ListRemove(kv *trib.KeyValue, n *int) error {
	return self.client.ListRemove(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, n)
}

func (self *keeperClient) ListKeys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix}
	if e := self.client.ListKeys(&np, list); e != nil {
		return e
	}
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
	}
	return nil
}

func (self *keeperClient) Clock(atLeast uint64, ret *uint64) error {
	return self.client.Clock(atLeast, ret)
}

func (self *keeperClient) Init() {
	sort.Sort(ByHashValue(self.backs))
}

func (self *keeperClient) KeeperBin(name string) []trib.Storage {
	self.Init()
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
	result := make([]trib.Storage, 3)
	for count < 3 {
		if result[count] = self.checkAddr(name, self.backs[startIndex]); result[count] != nil {
			count++
		}
		startIndex++
		if startIndex >= size {
			startIndex = 0
		}
	}
	return result
}

func (self *keeperClient) checkAddr(binName string, addr string) trib.Storage {
	client := &attClient{bin: binName, client: NewClient(addr)}
	list := trib.List{[]string{}}
	err := client.Keys(&trib.Pattern{"", "Completed"}, &list)
	if err == nil {
		return client
	} else {
		return nil
	}
}
