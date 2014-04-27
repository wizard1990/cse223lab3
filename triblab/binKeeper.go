package triblab

import (
    "trib"
    "fmt"
    "time"
)

type binKeeper struct {
    backs []string
    Ready chan<- bool
}

func NewKeeper(kc *trib.KeeperConfig) *binKeeper {
    keeper := binKeeper{backs:make([]string, len(kc.Backs)), Ready:kc.Ready}
    copy(keeper.backs, kc.Backs)
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
        err := <- testChan
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
            case <- timer:
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