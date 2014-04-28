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

func (self *binKeeper) run() error {
	testChan := make(chan error, len(self.backs))
	bc := NewBinClient(self.backs)

	//check the connection to each back-end
	for _, addr := range self.backs {
		go func(addr string) {
			var t uint64 = 0
			testChan <- bc.Bin(addr).Clock(0, &t)
		}(addr)
	}
	for i := 0; i < len(self.backs); i++ {
		err := <-testChan
		if err != nil {
			if self.Ready != nil {
				self.Ready <- false
			}
			return err
		}
	}
	if self.Ready != nil {
		self.Ready <- true
	}

	//start keeping clock
	var synClock uint64 = 0
	timer := time.Tick(1000 * time.Millisecond)
	results := make(chan uint64, len(self.backs))

	for {
		select {
		case <-timer:
			for _, addr := range self.backs {
				go func(addr string) {
					var t uint64 = 0
					bc.Bin(addr).Clock(synClock, &t)
					results <- t
				}(addr)
			}
			go func() {
				for i := 0; i < len(self.backs); i++ {
					t := <-results
					if t > synClock {
						synClock = t
					}
				}
			}()
		default:
			time.Sleep(100 * time.Millisecond)
		}

	}
	return fmt.Errorf("Warning! Big brother is not watching you!")
}
