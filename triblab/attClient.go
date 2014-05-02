package triblab

import (
	"trib"
	"trib/colon"
)

type attClient struct {
	bin 		string
	binManager	*binHash
	client 		[]trib.Storage
	channel 	chan int
}

func genPrefix(s string) string {
	return colon.Escape(s) + "::"
}

func (self *attClient) RefreshBin() {
	if self.binManager != nil {
		self.client = self.binManager.GetBinCopies(self.bin)
		self.channel = make(chan int, 3)
	}
}

func (self *attClient) Get(key string, value *string) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}

	go func(c chan int) {
		e1 := self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::KV", &res0)
		if e1 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e2 := self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::KV", &res1)
		if e2 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e3 := self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::KV", &res2)
		if e3 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	s := 0
	for i := 0; i < 3; i++ {
		s += (<-self.channel)
	}

	_, _, ele := FindLargestClock(&res0, &res1, &res2)

	*value = ele
	if s < 3 {
		self.RefreshBin()
	}
	return nil
}

func (self *attClient) Set(kv *trib.KeyValue, succ *bool) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	go func(c chan int){ 
		self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::KV", &res0)
		c <- 1
	}(self.channel)
	go func(c chan int){ 
		self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::KV", &res1)
		c <- 1
	}(self.channel)
	go func(c chan int){ 
		self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::KV", &res2)
		c <- 1
	}(self.channel)
	for i := 0; i < 3; i++ {
		<-self.channel
	}
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
		go func(c chan int, i int) {
			e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin)+colon.Escape(kv.Key) + "::KV", AddClock(n, kv.Value)}, succ)
			c <- 1
			if e != nil {
				flag = true
			}
		}(self.channel, i)
	}
	for i := 0; i < 3; i++ {
		<-self.channel
	}
	if flag {
		self.RefreshBin()
	}
	*succ = true

	return nil
}

func (self *attClient) Keys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + colon.Escape(p.Prefix), colon.Escape(p.Suffix) + "::KV"}
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}

	go func(c chan int) {
		e1 := self.client[0].ListKeys(&np, &res0)
		if e1 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e2 := self.client[1].ListKeys(&np, &res1)
		if e2 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e3 := self.client[2].ListKeys(&np, &res2)
		if e3 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	s := 0
	for i := 0; i < 3; i++ {
		s += (<-self.channel)
	}

	if s < 3 {
		self.RefreshBin()
	}
	list.L = (MergeKeyList(&res0, &res1, &res2)).L
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
		list.L[i] = colon.Unescape(list.L[i][:len(list.L[i])-4])
	}
	return nil
}

func (self *attClient) ListGet(key string, list *trib.List) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	go func(c chan int) {
		e1 := self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::L", &res0)
		if e1 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e2 := self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::L", &res1)
		if e2 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e3 := self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(key)+"::L", &res2)
		if e3 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	s := 0
	for i := 0; i < 3; i++ {
		s += (<-self.channel)
	}
	if s < 3 {
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
	go func(c chan int) {
		self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res0)
		c <- 1
	}(self.channel)
	go func(c chan int) {
		self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res1)
		c <- 1
	}(self.channel)
	go func(c chan int) {
		self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res2)
		c <- 1
	}(self.channel)
	for i := 0; i < 3; i++ {
		<-self.channel
	}
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
		go func(c chan int, i int) {
			e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + colon.Escape(kv.Key) + "::L", AddClock(n, kv.Value) + "::Append"}, succ)
			c <- 1
			if e != nil {
			flag = true
		}
		}(self.channel, i)
	}
	for i := 0; i < 3; i++ {
		<-self.channel
	}
	if flag {
		self.RefreshBin()
	}
	*succ = true
	return nil
}

func (self *attClient) ListRemove(kv *trib.KeyValue, n *int) error {
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	go func(c chan int) {
		self.client[0].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res0)
		c <- 1
	}(self.channel)
	go func(c chan int) {
		self.client[1].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res1)
		c <- 1
	}(self.channel)
	go func(c chan int) {
		self.client[2].ListGet(genPrefix(self.bin)+colon.Escape(kv.Key)+"::L", &res2)
		c <- 1
	}(self.channel)
	for i := 0; i < 3; i++ {
		<-self.channel
	}
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
		go func(c chan int, i int) {
			e := self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + colon.Escape(kv.Key) + "::L", AddClock(maxClk, kv.Value) + "::Remove"}, &succ)
			if e != nil {
				flag = true
			}
			c <- 1
		}(self.channel, i)
	}
	for i := 0; i < 3; i++ {
		<-self.channel
	}
	if flag {
		self.RefreshBin()
	}
	return nil
}

func (self *attClient) ListKeys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + colon.Escape(p.Prefix), colon.Escape(p.Suffix) + "::L"}
	res0 := trib.List{[]string{}}
	res1 := trib.List{[]string{}}
	res2 := trib.List{[]string{}}
	go func(c chan int) {
		e1 := self.client[0].ListKeys(&np, &res0)
		if e1 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e2 := self.client[1].ListKeys(&np, &res1)
		if e2 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	go func(c chan int) {
		e3 := self.client[2].ListKeys(&np, &res2)
		if e3 == nil {
			c <- 1
		} else {
			c <- 0
		}
	}(self.channel)
	s := 0
	for i := 0; i < 3; i++ {
		s += (<-self.channel)
	}
	if s < 3 {
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
		go func(c chan int, i int) {
			e := self.client[i].Clock(atLeast, ret)
			if e != nil {
				flag = true
			}
			c <- 1
		}(self.channel, i)
	}
	if flag {
		self.RefreshBin()
	}
	return nil
}
