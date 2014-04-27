package triblab

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"time"

	//"sync"

//"trib"
)

type BoardcastBin interface {
	Ask(bin_name string, status *int) error
}

// RPC inside holder
type Keep_bin struct {
	keeper *binKeeper
	//bin_map map[string]int
}

//var _ BoardcastBin = new(Keep_bin)

//Get RPC function
func Get_RPC(b *binKeeper) BoardcastBin {
	return &Keep_bin{keeper: b}
}

func (self *Keep_bin) Ask(bin_name string, status *int) error {
	//fmt.Println("Asking for ", bin_name)
	value, ok := self.keeper.locked_bin[bin_name]
	if ok {
		*status = value
	} else {
		*status = 0
	}
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
	e = http.Serve(s, sv)
	if e != nil {
		self.Ready <- false
	}
	return e
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

func (self *binKeeper) start_audit_bin(bin_name string) bool {
	self.bin_lock.Lock()
	//Go look local map
	value, ok := self.locked_bin[bin_name]
	if ok && value != 0 {
		self.bin_lock.Unlock()
		return false
	}

	self.locked_bin[bin_name] = 2
	self.bin_lock.Unlock()

	status_list := make([]int, len(self.Keeper_addrs))

	//Start to boardcast
	for i, _ := range status_list {
		status_list[i] = 0
		self.Ask(self.Keeper_addrs[i], bin_name, &status_list[i])
	}

	locked_by_others := false
	for i, _ := range status_list {
		if status_list[i] != 0 {
			locked_by_others = true
			break
		}
	}

	if locked_by_others {
		self.bin_lock.Lock()
		self.locked_bin[bin_name] = 0
		self.bin_lock.Unlock()
		return false
	}

	self.bin_lock.Lock()
	self.locked_bin[bin_name] = 1
	self.bin_lock.Unlock()
	return true
}

func (self *binKeeper) end_audit_bin(bin_name string) bool {
	self.bin_lock.Lock()
	self.locked_bin[bin_name] = 0
	self.bin_lock.Unlock()
	return true
}

func (self *binKeeper) Tao_T() error {
	time.Sleep(1 * time.Second)
	/*
					   0 hold the "a"&"c" for 2 seconds
					   1 try a in 1 second -> fail
					   2 try a in 3 second -> succ
				       3 hold the "b" for 2 seconds
				       4 try b in 1 second -> fail and retry in 3 second -> succ
		               5 try to get "c" in 1 second -> fail
	*/
	if self.This == 0 {
		if self.start_audit_bin("a") {
			fmt.Println(self.This, "Lock a")
		} else {
			fmt.Println(self.This, "Fail to get a")
		}
		if self.start_audit_bin("c") {
			fmt.Println(self.This, "Lock c")
		} else {
			fmt.Println(self.This, "Fail to get c")
		}

		time.Sleep(2 * time.Second)
		self.end_audit_bin("a")
		self.end_audit_bin("c")
		fmt.Println(self.This, "Relese a")
	}
	if self.This == 1 {
		time.Sleep(1 * time.Second)
		if self.start_audit_bin("a") {
			fmt.Println(self.This, "Error to get a")
		} else {
			fmt.Println(self.This, "Fail to get a")
		}
		self.end_audit_bin("a")
	}
	if self.This == 2 {
		time.Sleep(3 * time.Second)
		if self.start_audit_bin("a") {
			fmt.Println(self.This, "Lock a")
			time.Sleep(2 * time.Second)
		} else {
			fmt.Println(self.This, "Error Fail to get a")
		}
		self.end_audit_bin("a")
		fmt.Println(self.This, "Relese a")
	}
	if self.This == 3 {
		if self.start_audit_bin("b") {
			fmt.Println(self.This, "Lock b")
			time.Sleep(2 * time.Second)
		} else {
			fmt.Println(self.This, "Error Fail to get b")
		}
		self.end_audit_bin("b")
		fmt.Println(self.This, "Relese b")
	}
	if self.This == 4 {
		time.Sleep(1 * time.Second)
		if self.start_audit_bin("b") {
			fmt.Println(self.This, "Error Lock b")
		} else {
			fmt.Println(self.This, "Fail to get b")
		}
		time.Sleep(2 * time.Second)
		if self.start_audit_bin("b") {
			fmt.Println(self.This, "Lock b")
		} else {
			fmt.Println(self.This, "Error Fail to get b")
		}
		fmt.Println(self.This, "Relese b")
	}
	if self.This == 5 {
		time.Sleep(1 * time.Second)
		if self.start_audit_bin("c") {
			fmt.Println(self.This, "Error Lock c")
		} else {
			fmt.Println(self.This, "Fail to get c")
		}

	}
	if self.This == 1 {
	}
	return nil
}
