package triblab

import (
	//"hash/fnv"
	"sort"
	"trib"
	//"trib/colon"
)

type keeperClient struct {
	backs    []string
	indexMap map[string]int
}

/*
func genPrefix(s string) string {
	return colon.Escape(s) + "::"
}
*/

type KeeperStorage interface {
	KeeperBin(name string) ([]trib.Storage, []string)
}

/*
 */
func (self *keeperClient) Init() {
	sort.Sort(ByHashValue(self.backs))
}

func (self *keeperClient) KeeperBin(name string) ([]trib.Storage, []string) {
	self.Init()
	if self.backs == nil {
		return nil,nil
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
	addrResult := make([]string, 3)
	for count < 3 {
		if result[count] = self.checkAddr(name, self.backs[startIndex]); result[count] != nil {
			addrResult[count] = self.backs[startIndex]
			count++
		}
		startIndex++
		if startIndex >= size {
			startIndex = 0
		}
	}
	return result, addrResult
}

func (self *keeperClient) checkAddr(binName string, addr string) trib.Storage {
	client := &attKeeperClient{bin: binName, client: NewClient(addr)}
	list := trib.List{[]string{}}
	err := client.Keys(&trib.Pattern{"", "Completed"}, &list)
	if err == nil {
		return client
	} else {
		return nil
	}
}
