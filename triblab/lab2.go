package triblab

import (
	"sync"
	"trib"
)

func NewBinClient(backs []string) trib.BinStorage {
	return &binClient{backs: backs, indexMap: make(map[string]trib.Storage)}
}

func ServeKeeper(kc *trib.KeeperConfig) error {
	keeper := NewKeeper(kc)
	go keeper.Tao_T()
	e := keeper.Serve_Consistency_RPC()
	if e != nil {
		return e
	}
	return keeper.run()
}

func NewFront(s trib.BinStorage) trib.Server {
	return &binServer{server: s, lock: sync.Mutex{}}
}
