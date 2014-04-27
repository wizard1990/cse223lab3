package triblab

import (
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"trib"
)

// Creates an RPC client that connects to addr.
func NewClient(addr string) trib.Storage {
	return &client{addr: addr, lock: sync.Mutex{}}
}

// Serve as a backend based on the given configuration
func ServeBack(b *trib.BackConfig) error {
	sv := rpc.NewServer()
	e := sv.Register(b.Store)
	if e != nil {
		go ReadyNotify(b.Ready, false)
		return e
	}
	s, e := net.Listen("tcp", b.Addr)
	if e != nil {
		go ReadyNotify(b.Ready, false)
		return e
	}
	go ReadyNotify(b.Ready, true)
	return http.Serve(s, sv)
}

func ReadyNotify(c chan<- bool, b bool) {
	if c != nil {
		c <- b
	}
}
