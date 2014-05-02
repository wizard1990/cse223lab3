package triblab

import (
	"net/rpc"
	"sync"
	"trib"
)

var CONNECT_TRIALS int = 0

type client struct {
	addr string
	conn *rpc.Client
	lock sync.Mutex
}

func (self *client) Connect(reconnect bool) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if (self.conn != nil) && (!reconnect) {
		return nil
	}
	var e error
	for i := 0; i <= CONNECT_TRIALS; i++ {
		self.conn, e = rpc.DialHTTP("tcp", self.addr)
		if e == nil {
			break
		}
	}
	return e
}

func (self *client) Get(key string, value *string) error {
	if e := self.Connect(false); e != nil {
		//fmt.Println(e)
		return e
	}

	*value = ""
	count := 0
	for e := self.conn.Call("Storage.Get", key, value); e != nil; {
		if count < CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	return nil
}

func (self *client) Set(kv *trib.KeyValue, succ *bool) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	count := 0
	for e := self.conn.Call("Storage.Set", kv, succ); e != nil; {
		if count < CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	return nil
}

func (self *client) Keys(p *trib.Pattern, list *trib.List) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	var tmpList *trib.List = &trib.List{[]string{}}
	count := 0
	for e := self.conn.Call("Storage.Keys", p, tmpList); e != nil; {
		if count < CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	*list = *tmpList
	self.conn.Close()
	return nil
}

func (self *client) ListGet(key string, list *trib.List) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	list.L = nil
	count := 0
	for e := self.conn.Call("Storage.ListGet", key, list); e != nil; {
		if count < CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	if list.L == nil {
		list.L = make([]string, 0)
	}
	self.conn.Close()
	return nil
}

func (self *client) ListAppend(kv *trib.KeyValue, succ *bool) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	count := 0
	for e := self.conn.Call("Storage.ListAppend", kv, succ); e != nil; {
		if count < CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	self.conn.Close()
	return nil
}

func (self *client) ListRemove(kv *trib.KeyValue, n *int) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	count := 0
	for e := self.conn.Call("Storage.ListRemove", kv, n); e != nil; {
		if count < CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	self.conn.Close()
	return nil
}

func (self *client) ListKeys(p *trib.Pattern, list *trib.List) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	list.L = nil
	count := 0
	for e := self.conn.Call("Storage.ListKeys", p, list); e != nil; {
		if count > CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	if list.L == nil {
		list.L = make([]string, 0)
	}
	self.conn.Close()
	return nil
}

func (self *client) Clock(atLeast uint64, ret *uint64) error {
	if e := self.Connect(true); e != nil {
		return e
	}
	count := 0
	for e := self.conn.Call("Storage.Clock", atLeast, ret); e != nil; {
		if count > CONNECT_TRIALS {
			return e
		}
		e = self.Connect(true)
		if e != nil {
			self.conn.Close()
			self.conn = nil
			return e
		}
		count++
	}
	self.conn.Close()
	return nil
}
