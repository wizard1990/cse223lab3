package entries

import (
	. "trib"

	"triblab"
)

// Makes a front end that talks to one single backend
// Used in Lab1
func MakeFrontSingle(back string) Server {
	return triblab.NewFront(triblab.NewBinClient([]string{back}))
}

// Serve as a single backend.
// Listen on addr, using s as underlying storage.
func ServeBackSingle(addr string, s Storage, ready chan<- bool) error {
	back := &BackConfig{
		Addr:  addr,
		Store: s,
		Ready: ready,
	}

	return triblab.ServeBack(back)
}
