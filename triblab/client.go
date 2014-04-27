package triblab

import (
    "fmt"
    "trib"
    "net/rpc"
    "hash/fnv"
    "sync"
)

var CONNECT_TRIALS int = 5

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
    for i := 0; i < CONNECT_TRIALS; i++ {
        self.conn, e = rpc.DialHTTP("tcp", self.addr)
        if e == nil {
            break
        }
    }
    if e != nil {
        fmt.Println(e)
    }
    return e
}

func (self *client) Get(key string, value *string) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }

    *value = ""
    for e := self.conn.Call("Storage.Get", key, value); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    return nil
}

func (self *client) Set(kv *trib.KeyValue, succ *bool) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    for e := self.conn.Call("Storage.Set", kv, succ); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    return nil
}

func (self *client) Keys(p *trib.Pattern, list *trib.List) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    var tmpList *trib.List = &trib.List{[]string{}}
    for e := self.conn.Call("Storage.Keys", p, tmpList); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    *list = *tmpList
    return nil
}

func (self *client) ListGet(key string, list *trib.List) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    var tmpList *trib.List = &trib.List{[]string{}} 
    for e := self.conn.Call("Storage.ListGet", key, tmpList); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    *list = *tmpList
    return nil
}

func (self *client) ListAppend(kv *trib.KeyValue, succ *bool) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    for e := self.conn.Call("Storage.ListAppend", kv, succ); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    return nil
}

func (self *client) ListRemove(kv *trib.KeyValue, n *int) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    for e := self.conn.Call("Storage.ListRemove", kv, n); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    return nil
}

func (self *client) ListKeys(p *trib.Pattern, list *trib.List) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    for e := self.conn.Call("Storage.ListKeys", p, list); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    return nil
}

func (self *client) Clock(atLeast uint64, ret *uint64) error {
    if e := self.Connect(false); e != nil {
        fmt.Println(e)
        return e
    }
    for e := self.conn.Call("Storage.Clock", atLeast, ret); e != nil; {
        e = self.Connect(true)
        if e != nil {
            return e
        }
    }
    return nil
}
