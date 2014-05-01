package triblab

import (
	"fmt"
	"trib"
	"trib/colon"
)

type attClient struct {
	bin 		string
	binManager	*binHash
	client 		[]trib.Storage
}

func genPrefix(s string) string {
	return colon.Escape(s) + "::"
}

func (self *attClient) RefreshBin() {
	if self.binManager != nil {
		self.client = self.binManager.GetBinCopies(self.bin)
	}
}

func (self *attClient) Get(key string, value *string) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	e1 := self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::KV", &res0)
	e2 := self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::KV", &res1)
	e3 := self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::KV", &res2)
	_, _, ele := FindLargestClock(&res0, &res1, &res2)

	*value = ele
	if (e1 != nil) || (e2 != nil) || (e3 != nil) {
		self.RefreshBin()
	}
	return nil
}

func (self *attClient) Set(kv *trib.KeyValue, succ *bool) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::KV", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::KV", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::KV", &res2)
	clk, _, _ := FindLargestClock(&res0, &res1, &res2)

	clk++ //YAN remeber to increase the clk before call Lock
	var n uint64
	for i := 0; i < 3; i++ {
		var t uint64
		self.client[i].Clock(clk, &t)
		if t > n {
			n = t
		}
	}
	flag := false
	for i := 0; i < 3; i++ {
		e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin)+colon.Escape(kv.Key) + "::KV", AddClock(n, kv.Value)}, succ)
		if e != nil {
			flag = true
		}
	}
	if flag {
		self.RefreshBin()
	}

	return nil
}

func (self *attClient) Keys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + colon.Escape(p.Prefix), p.Suffix + "::KV"}
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	e1 := self.client[0].ListKeys(&np, &res0)
	e2 := self.client[1].ListKeys(&np, &res1)
	e3 := self.client[2].ListKeys(&np, &res2)
	if (e1 != nil) || (e2 != nil) || (e3 != nil) {
		self.RefreshBin()
	}
	list.L = (MergeKeyList(&res0, &res1, &res2)).L
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
		list.L[i] = colon.Unescape(list.L[i][:len(list.L[i])-4])
		//fmt.Println(list.L)
	}
	return nil
}

func (self *attClient) ListGet(key string, list *trib.List) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	e1 := self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::L", &res0)
	e2 := self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::L", &res1)
	e3 := self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::L", &res2)
	if (e1 != nil) || (e2 != nil) || (e3 != nil) {
		self.RefreshBin()
	}
	_, res, _ := FindLargestClock(&res0, &res1, &res2)

	list.L = GetDisplayList(res).L

	return nil
}

func (self *attClient) ListAppend(kv *trib.KeyValue, succ *bool) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res2)
	clk, _, _ := FindLargestClock(&res0, &res1, &res2)
	clk++

	var n uint64
	for i := 0; i < 3; i++ {
		var t uint64
		self.client[i].Clock(clk, &t)
		if t > n {
			n = t
		}
	}
	flag := false
	for i := 0; i < 3; i++ {
		e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + colon.Escape(kv.Key) + "::L", AddClock(n, kv.Value) + "::Append"}, succ)
		if e != nil {
			flag = true
		}
	}
	if flag {
		self.RefreshBin()
	}
	return nil
}

func (self *attClient) ListRemove(kv *trib.KeyValue, n *int) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res0)
	self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res1)
	self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res2)
	clk, res, _ := FindLargestClock(&res0, &res1, &res2)
	clk++

	var maxClk uint64
	for i := 0; i < 3; i++ {
		var t uint64
		self.client[i].Clock(clk, &t)
		if t > maxClk {
			maxClk = t
		}
	}

	resCnt := 0
	t := 0
	for _, v := range res.L {
		if _, tv := SplitClock(v); tv[:len(tv)-8] == kv.Value {
			for i := 0; i < 3; i++ {
				self.client[i].ListRemove(&trib.KeyValue{genPrefix(self.bin) + colon.Escape(kv.Key) + "::L", v}, &t)
			}
			if tv[len(tv)-6:] == "Append" {
				resCnt++
			} else {
				resCnt = 0
			}
		}
	}
	*n = resCnt
	flag := false
	for i := 0; i < 3; i++ {
		succ := false
		e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + colon.Escape(kv.Key) + "::L", AddClock(maxClk, kv.Value) + "::Remove"}, &succ)
		if e != nil {
			flag = true
		}
	}
	if flag {
		self.RefreshBin()
	}
	return nil
}

func (self *attClient) ListKeys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + colon.Escape(p.Prefix), p.Suffix + "::L"}
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	e1 := self.client[0].ListKeys(&np, &res0)
	e2 := self.client[1].ListKeys(&np, &res1)
	e3 := self.client[2].ListKeys(&np, &res2)
	if (e1 != nil) || (e2 != nil) || (e3 != nil) {
		self.RefreshBin()
	}

	list.L = (MergeKeyList(&res0, &res1, &res2)).L
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
		list.L[i] = colon.Unescape(list.L[i][:len(list.L[i])-3])
	}
	return nil
}

func (self *attClient) Clock(atLeast uint64, ret *uint64) error {
	flag := false
	for i := 0; i < 3; i++ {
		e := self.client[i].Clock(atLeast, ret)
		if e != nil {
			flag = true
		}
	}
	if flag {
		self.RefreshBin()
	}
	return nil
	fmt.Println("123")
	return nil
}
