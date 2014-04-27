package triblab

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

//"sync"
//"trib"
)

type Boardcast_bin interface {
	Ask(bin_name string, status *int) error
}

// RPC inside holder
type Keep_bin struct {
	hold_bin map[string]int
	name     string
}

func (self *Keep_bin) Ask(bin_name string, status *int) error {
	fmt.Println("T", self.name)
	*status = 0
	return nil
}

//Get RPC function
func Get_RPC() Boardcast_bin {
	return &Keep_bin{name: "123"}
}

func (self *binKeeper) Tao_T() error {
	fmt.Println(self.This, self.Keeper_addrs)
	var status int
	rpc1 := Get_RPC()
	rpc1.Ask("1", &status)
	return nil
}
func (self *binKeeper) Serve_Consistency_RPC() error {
	keepbiner := Keep_bin{name: "123"}
	sv := rpc.NewServer()
	e := sv.Register(keepbiner)
	if e != nil {
		self.Ready <- false
	}
	s, e := net.Listen("tcp", self.This_Addr)
	if e != nil {
		self.Ready <- false
	}
	return http.Serve(s, sv)
}
func (self *binKeeper) start_audit_bin(bin_name string) bool {
	return true
}

func (self *binKeeper) end_audit_bin(bin_name string) bool {
	return true
}
