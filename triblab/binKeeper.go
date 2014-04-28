package triblab

import (
	"fmt"
	"sync"
	"time"
	"trib"
)

type binKeeper struct {
	backs []string
	Ready chan<- bool

	clock_lock sync.Mutex
	MaxCount   uint64

	//Args for consistency
	bin_lock   sync.Mutex
	locked_bin map[string]int
	//End consisntecy

	//Xintian Args for replicate
	clientMap map[string]trib.Storage
	bc        KeeperStorage
	//End Xintian

	Keeper_addrs []string // Keepers peers' addr not included myself
	This         int      //My index // useless
	This_Addr    string
}

func NewKeeper(kc *trib.KeeperConfig) *binKeeper {
	keeper := binKeeper{backs: make([]string, len(kc.Backs)), Ready: kc.Ready}
	keeper.Keeper_addrs = make([]string, len(kc.Addrs)-1)
	keeper.This = kc.This
	keeper.locked_bin = make(map[string]int)

	j := 0
	for i, _ := range kc.Addrs {
		if i == kc.This {
			keeper.This_Addr = kc.Addrs[i]
			continue
		}
		keeper.Keeper_addrs[j] = kc.Addrs[i]
		j++
	}
	copy(keeper.backs, kc.Backs)

	//Xintian for replicate
	//Cancel by Tao
	//Plz don't commit error code

	/*
			keeper.bc = NewKeeperClient(kc.Backs) // keeper client
			for _,addr := range keeper.backs{
		    keeper.clientMap[addr] = &client{addr:addr}
			}
	*/
	//End Xintian

	return &keeper
}

func (self *binKeeper) send_clock_sync(addr string, c uint64) error {
	rpc := NewClient(addr)
	var cret uint64
	fmt.Println(addr, c)
	e := rpc.Clock(c, &cret)

	self.clock_lock.Lock()
	if cret > self.MaxCount {
		self.MaxCount = cret
	}
	self.clock_lock.Unlock()

	return e
}
func (self *binKeeper) clock_sync() error {
	fmt.Println("clock_sync")

	for {
		for _, addr := range self.backs {
			go self.send_clock_sync(
				addr, self.MaxCount)
		}
		time.Sleep(1 * time.Second)
	}
	return nil

}
