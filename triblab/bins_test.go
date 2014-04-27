package triblab_test

import (
	"testing"
	"fmt"

	"triblab"
	"trib/entries"
	"trib/randaddr"
	"trib/store"
	"trib/tribtest"
	"trib"
)

func TestBinStorage(t *testing.T) {
	addr1 := randaddr.Local()
	addr2 := randaddr.Local()

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
