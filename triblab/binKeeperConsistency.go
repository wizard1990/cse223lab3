package triblab

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"time"

	"sync"

//"trib"
)

type BoardcastBin interface {
	Ask(bin_name string, status *int) error
}

// RPC inside holder
type Keep_bin struct {
	keeper *binKeeper
	lock   sync.Mutex
	bin    map[string]int
}

//var _ BoardcastBin = new(Keep_bin)

//Get RPC function
func Get_RPC(b *binKeeper) BoardcastBin {
	bin_map := make(map[string]int)
	return &Keep_bin{keeper: b, bin: bin_map}
}

func (self *Keep_bin) Ask(bin_name string, status *int) error {
	fmt.Println("Asking for ", bin_name)
	*status = 1
	return nil
}
func (self *binKeeper) Serve_Consistency_RPC() error {
	keepbiner := Get_RPC(self)
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

func (self *binKeeper) Ask(server, bin_name string, status *int) error {
	conn, e := rpc.DialHTTP("tcp", server)
	if e != nil {
		return e
	}
	e = conn.Call("Keep_bin.Ask", bin_name, status)
	if e != nil {
		return e
	}
	return nil
}
func (self *binKeeper) Tao_T() error {

	time.Sleep(1 * time.Second)
	var status int
	for _, addr := range self.Keeper_addrs {
		fmt.Println(addr)
		self.Ask(addr, "a", &status)
		fmt.Println(status)
	}
	return nil
}
func (self *binKeeper) start_audit_bin(bin_name string) bool {
	return true
}

func (self *binKeeper) end_audit_bin(bin_name string) bool {
	return true
}
