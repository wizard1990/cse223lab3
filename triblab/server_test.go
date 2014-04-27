package triblab_test

import (
	"testing"

	"trib/entries"
	"trib/randaddr"
	"trib/store"
	"trib/tribtest"
)

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
	server := entries.MakeFrontSingle(addr)

	tribtest.CheckServerConcur(t, server)
}
