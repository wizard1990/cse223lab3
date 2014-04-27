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
