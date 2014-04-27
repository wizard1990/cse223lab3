package triblab

import (
	"trib"
    "fmt"
    "sync"
)

func NewBinClient(backs []string) trib.BinStorage {
	return &binClient{backs: backs}
}

func ServeKeeper(kc *trib.KeeperConfig) error {
	keeper := NewKeeper(kc)
    return keeper.run()
}

func NewFront(s trib.BinStorage) trib.Server {
	return &binServer{server:s, lock:sync.Mutex{}}
}