package triblab_test

import (
	"fmt"
	"testing"

	"trib"
	"trib/entries"
	"trib/randaddr"
	"trib/store"
	"trib/tribtest"
	"triblab"
)

func TestBinStorage(t *testing.T) {
	addr1 := randaddr.Local()
	addr2 := randaddr.Local()
	for addr2 == addr1 {
		addr2 = randaddr.Local()
	}

	ready1 := make(chan bool)
	ready2 := make(chan bool)

	run := func(addr string, ready chan bool) {
		e := entries.ServeBackSingle(addr,
			store.NewStorage(), ready)
		if e != nil {
			t.Fatal(e)
		}
	}

	go run(addr1, ready1)
	go run(addr2, ready2)

	r := <-ready1 && <-ready2
	if !r {
		t.Fatal("not ready")
	}

	readyk := make(chan bool)
	addrk := randaddr.Local()
	for addrk == addr1 || addrk == addr2 {
		addrk = randaddr.Local()
	}
	go func() {
		e := triblab.ServeKeeper(&trib.KeeperConfig{
			Backs: []string{addr1, addr2},
			Addrs: []string{addrk},
			This:  0,
			Id:    0,
			Ready: readyk,
		})
		if e != nil {
			t.Fatal(e)
		}
	}()

	if !<-readyk {
		t.Fatal("keeper not ready")
	}

	bc := triblab.NewBinClient(
		[]string{addr1, addr2},
	)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		c := bc.Bin(fmt.Sprintf("b%d", i))
		go func(s trib.Storage) {
			tribtest.CheckStorage(t, s)
			done <- true
		}(c)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
