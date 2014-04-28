package triblab

import (
    "trib"
    "trib/colon"
)

type attClient struct {
    bin string
    client []trib.Storage
}

func genPrefix(s string) string {
    return colon.Escape(s) + "::"
}

func (self *attClient) Get(key string, value *string) error {
    res := ""
    index := 0
    for i := 0; i < 3; i++ {
        self.client[i].Get(genPrefix(self.bin) + "Completed", &res)
        if len(res) == 1 {
            index = i
            break
        }
    }
    return self.client[index].Get(genPrefix(self.bin) + key, value)
}

func (self *attClient) Set(kv *trib.KeyValue, succ *bool) error {
    for i := 0; i < 3; i++ {
        self.client[i].Set(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, succ)
    }
    return nil
}

func (self *attClient) Keys(p *trib.Pattern, list *trib.List) error {
    res := ""
    index := 0
    for i := 0; i < 3; i++ {
        self.client[i].Get(genPrefix(self.bin) + "Completed", &res)
        if len(res) == 1 {
            index = i
            break
        }
    }
    np := trib.Pattern{genPrefix(self.bin)+p.Prefix, p.Suffix}
    if e := self.client[index].Keys(&np, list); e != nil {
        return e
    }
    for i, s := range list.L {
        list.L[i] = s[len(genPrefix(self.bin)):]
    }
    return nil
}

func (self *attClient) ListGet(key string, list *trib.List) error {
    res := ""
    index := 0
    for i := 0; i < 3; i++ {
        self.client[i].Get(genPrefix(self.bin) + "Completed", &res)
        if len(res) == 1 {
            index = i
            break
        }
    }
    return self.client[index].ListGet(genPrefix(self.bin) + key, list)
}

func (self *attClient) ListAppend(kv *trib.KeyValue, succ *bool) error {
    for i := 0; i < 3; i++ {
        self.client[i].ListAppend(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, succ)
    }
    return nil
}

func (self *attClient) ListRemove(kv *trib.KeyValue, n *int) error {
    for i := 0; i < 3; i++ {
        self.client[i].ListRemove(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, n)
    }
    return nil
}

func (self *attClient) ListKeys(p *trib.Pattern, list *trib.List) error {
    res := ""
    index := 0
    for i := 0; i < 3; i++ {
        self.client[i].Get(genPrefix(self.bin) + "Completed", &res)
        if len(res) == 1 {
            index = i
            break
        }
    }
    np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix}
    if e := self.client[index].ListKeys(&np, list); e != nil {
        return e
    }
    for i, s := range list.L {
        list.L[i] = s[len(genPrefix(self.bin)):]
    }
    return nil
}

func (self *attClient) Clock(atLeast uint64, ret *uint64) error {
    for i := 0; i < 3; i++ {
        return self.client[i].Clock(atLeast, ret)
    }
    return nil
}
