package triblab

import (
	"fmt"
	"sync"
	"trib"
)

func NewBinClient(backs []string) trib.BinStorage {
	return &binClient{backs: backs}
}

func NewKeeperClient(backs []string) KeeperStorage {
	return &keeperClient{backs: backs, indexMap: make(map[string]int)}
}

func ServeKeeper(kc *trib.KeeperConfig) error {
	keeper := NewKeeper(kc)
	//go keeper.Tao_T() //Test code for consistency
	go keeper.clock_sync()
	fmt.Println("servekeeper")
	go keeper.Replicate_bin()
	kc.Ready <- true
	return keeper.Serve_Consistency_RPC()
}

func NewFront(s trib.BinStorage) trib.Server {
	return &binServer{server: s, lock: sync.Mutex{}}
}
