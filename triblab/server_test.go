package triblab_test

import (
	"testing"

	"trib"
	"trib/entries"
	"trib/randaddr"
	"trib/store"
	"trib/tribtest"
	"triblab"
)

func startKeeper(t *testing.T, addr string) {
	readyk := make(chan bool)
	addrk := randaddr.Local()
	for addrk == addr {
		addrk = randaddr.Local()
	}

	go func() {
		e := triblab.ServeKeeper(&trib.KeeperConfig{
			Backs: []string{addr},
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
}

func TestServer(t *testing.T) {
	addr := randaddr.Local()
	ready := make(chan bool)
	go func() {
		e := entries.ServeBackSingle(addr, store.NewStorage(), ready)
		if e != nil {
			t.Fatal(e)
		}
	}()
	<-ready

	startKeeper(t, addr)

	server := entries.MakeFrontSingle(addr)

	tribtest.CheckServer(t, server)
}

func TestServerConcur(t *testing.T) {
	addr := randaddr.Local()
	ready := make(chan bool)
	go func() {
		e := entries.ServeBackSingle(addr, store.NewStorage(), ready)
		if e != nil {
			t.Fatal(e)
		}
	}()

	<-ready

	startKeeper(t, addr)

	server := entries.MakeFrontSingle(addr)

	tribtest.CheckServerConcur(t, server)
}
