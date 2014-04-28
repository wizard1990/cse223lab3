package triblab

import (
	"trib"
)

type attKeeperClient struct {
	bin    string
	client trib.Storage
}

func (self *attKeeperClient) Get(key string, value *string) error {
	return self.client.Get(genPrefix(self.bin)+key, value)
}

func (self *attKeeperClient) Set(kv *trib.KeyValue, succ *bool) error {
	return self.client.Set(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, succ)
}

func (self *attKeeperClient) Keys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix}
	if e := self.client.Keys(&np, list); e != nil {
		return e
	}
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
	}
	return nil
}

func (self *attKeeperClient) ListGet(key string, list *trib.List) error {
	return self.client.ListGet(genPrefix(self.bin)+key, list)
}

func (self *attKeeperClient) ListAppend(kv *trib.KeyValue, succ *bool) error {
	return self.client.ListAppend(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, succ)
}

func (self *attKeeperClient) ListRemove(kv *trib.KeyValue, n *int) error {
	return self.client.ListRemove(&trib.KeyValue{genPrefix(self.bin) + kv.Key, kv.Value}, n)
}
func (self *attKeeperClient) ListKeys(p *trib.Pattern, list *trib.List) error {
	np := trib.Pattern{genPrefix(self.bin) + p.Prefix, p.Suffix}
	if e := self.client.ListKeys(&np, list); e != nil {
		return e
	}
	for i, s := range list.L {
		list.L[i] = s[len(genPrefix(self.bin)):]
	}
	return nil
}

func (self *attKeeperClient) Clock(atLeast uint64, ret *uint64) error {
	return self.client.Clock(atLeast, ret)
}
