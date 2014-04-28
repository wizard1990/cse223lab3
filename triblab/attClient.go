package triblab

import (
	"fmt"
	"trib"
	"trib/colon"
)

type attClient struct {
	bin    string
	client []trib.Storage
}

func genPrefix(s string) string {
	return colon.Escape(s) + "::"
}

func (self *attClient) Get(key string, value *string) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+key+"::KV", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+key+"::KV", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+key+"::KV", &res2)
	_, _, ele := FindLargestClock(&res0, &res1, &res2)

	*value = ele
	return nil
}

func (self *attClient) Set(kv *trib.KeyValue, succ *bool) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	e := self.client[0].ListGet(genPrefix(self.bin)+kv.Key+"::KV", &res0)
	fmt.Println(e)
	self.client[1].ListGet(genPrefix(self.bin)+kv.Key+"::KV", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+kv.Key+"::KV", &res2)
	clk, _, _ := FindLargestClock(&res0, &res1, &res2)

	var n uint64
	for i := 0; i < 3; i++ {
		var t uint64
		self.client[i].Clock(clk, &t)
		if t > n {
			n = t
		}
	}
	for i := 0; i < 3; i++ {
		succ := false
		e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + kv.Key + "::KV", AddClock(n, kv.Value)}, &succ)
		fmt.Println(e)
		fmt.Println(succ)
	}
	return nil
}

func (self *attClient) Keys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix + "::KV"}
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListKeys(&np, &res0)
	self.client[1].ListKeys(&np, &res1)
	self.client[2].ListKeys(&np, &res2)
	list.L = (MergeKeyList(&res0, &res1, &res2)).L
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
		list.L[i] = list.L[i][:len(list.L[i]) - 4]
		fmt.Println(list.L)
	}
	return nil
}

func (self *attClient) ListGet(key string, list *trib.List) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+key+"::L", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+key+"::L", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+key+"::L", &res2)
	_, res, _ := FindLargestClock(&res0, &res1, &res2)
	list.L = res.L
	for i, _ := range list.L {
		_, list.L[i] = SplitClock(list.L[i])
	}
	return nil
}

func (self *attClient) ListAppend(kv *trib.KeyValue, succ *bool) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+kv.Key+"::L", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+kv.Key+"::L", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+kv.Key+"::L", &res2)
	clk, _, _ := FindLargestClock(&res0, &res1, &res2)

	var n uint64
	for i := 0; i < 3; i++ {
		var t uint64
		self.client[i].Clock(clk, &t)
		if t > n {
			n = t
		}
	}
	for i := 0; i < 3; i++ {
		succ := false
		self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + kv.Key + "::L", AddClock(n, kv.Value)}, &succ)
	}
	return nil
}

func (self *attClient) ListRemove(kv *trib.KeyValue, n *int) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+kv.Key+"::L", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+kv.Key+"::L", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+kv.Key+"::L", &res2)
	_, res, _ := FindLargestClock(&res0, &res1, &res2)

	resCnt := 0
	t := 0
	for _, v := range res.L {
		if _, tv := SplitClock(v); tv == kv.Value {
			for i := 0; i < 3; i++ {
				self.client[0].ListRemove(&trib.KeyValue{genPrefix(self.bin) + kv.Key + "::L", v}, &t)
			}
			resCnt++
		}
	}
	*n = resCnt
	return nil
}

func (self *attClient) ListKeys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix + "::L"}
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListKeys(&np, &res0)
	self.client[1].ListKeys(&np, &res1)
	self.client[2].ListKeys(&np, &res2)

	list.L = (MergeKeyList(&res0, &res1, &res2)).L
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
		list.L[i] = list.L[i][:len(list.L[i]) - 3]
	}
	return nil
}

func (self *attClient) Clock(atLeast uint64, ret *uint64) error {
	for i := 0; i < 3; i++ {
		self.client[i].Clock(atLeast, ret)
	}
	return nil
}
