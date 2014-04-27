package triblab

import (
	"trib"
    "sync"
)

func NewBinClient(backs []string) trib.BinStorage {
	return &binClient{backs: backs, indexMap: make(map[string]int)}
}

func ServeKeeper(kc *trib.KeeperConfig) error {
	keeper := NewKeeper(kc)
	go keeper.Tao_T()
	go keeper.Serve_Consistency_RPC()
	return keeper.run()
}

func NewFront(s trib.BinStorage) trib.Server {
    return &binServer{server:s, lock:sync.Mutex{}}
}
