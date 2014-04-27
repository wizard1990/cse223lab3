package triblab_test

import (
	"testing"

	"trib/entries"
	"trib/randaddr"
	"trib/store"
	"trib/tribtest"
	"triblab"
)

func TestRPC(t *testing.T) {
	addr := randaddr.Local()
	ready := make(chan bool)

	go func() {
		e := entries.ServeBackSingle(addr, store.NewStorage(), ready)
		if e != nil {
			t.Fatal(e)
		}
	}()

	r := <-ready
	if !r {
		t.Fatal("not ready")
	}

	c := triblab.NewClient(addr)

	tribtest.CheckStorage(t, c)
}
