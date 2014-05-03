package triblab

import (
	"fmt"
	"strings"
	"time"
	"trib"
	"trib/colon"
)

func (self *binKeeper) Replicate_bin() error {

	for {
		for index := 0; index < len(self.backs); index++ {
			backend := self.clientMap[self.backs[index]]
			tmp := trib.List{[]string{}}

			e := backend.ListKeys(&(trib.Pattern{"", ""}), &tmp)
			if e != nil {
				continue
			}

			//TODO get bin_name

			bin_map := make(map[string]int)

			for _, str := range tmp.L {
				p := strings.Index(str, "::")
				bin_name := str[:p]

				bin_name = colon.Unescape(bin_name)

				_, ok := bin_map[bin_name]
				if ok {
					continue
				} else {
					bin_map[bin_name] = 1
				}
			}
			for bin_name, _ := range bin_map {
				e = self.audit_bin(bin_name, backend.addr)
				if e == nil {
					//remove local bin_name
				}
			}
		}
		//fmt.Println("Done Replication")
		time.Sleep(1 * time.Second)
	}
	time.Sleep(1 * time.Second)
	return fmt.Errorf("replication stops for strange reasons")
}

func (self *binKeeper) audit_bin(bin_name string, back string) error {
	flag := self.start_audit_bin(bin_name)
	defer self.end_audit_bin(bin_name)
	if flag == false {
		return fmt.Errorf("Fail to lock")
	}

	backends_client, backends_addrs := self.bc.KeeperBin(bin_name)

	_keys_list := make([]trib.List, 3)

	key_map := make(map[string]int)
	for i := 0; i < 3; i++ {
		e := backends_client[i].ListKeys(&trib.Pattern{}, &_keys_list[i])
		if e != nil {
			return e
		}
		for _, str := range _keys_list[i].L {
			_, ok := key_map[str]
			if ok {
				continue
			} else {
				key_map[str] = 1
			}
		}
	}
	all_replicate_done := true
	for key, _ := range key_map {
		e := self.replicate_key(key, backends_client)
		if e != nil {
			all_replicate_done = false
		}
	}

	if all_replicate_done {
		go self.remove_redundant(back, backends_addrs, bin_name)
		return nil
	} else {
		return fmt.Errorf("error")
	}

	fmt.Println(bin_name, backends_addrs)
	return nil
	fmt.Println(backends_client)
	return nil
}

func (self *binKeeper) replicate_key(key string, backends_client []trib.Storage) error {
	//Get all list
	lists := make([]trib.List, 3)
	for i := 0; i < 3; i++ {
		e := backends_client[i].ListGet(key, &lists[i])
		if e != nil {
			return e
		}
	}

	_, max_set, _ := FindLargestClock(&lists[0], &lists[1], &lists[2])

	for i := 0; i < 3; i++ {
		to_add := DiffList(max_set, &lists[i])

		for _, value := range to_add.L {
			succ := false
			e := backends_client[i].ListAppend(&trib.KeyValue{key, value}, &succ)
			if e != nil {
				return e
			}
			if succ == false {
				return fmt.Errorf("error0")
			}
		}
	}
	return nil
}

func (self *binKeeper) remove_redundant(back string, addrs []string, bin_name string) error {
	for i := 0; i < 3; i++ {
		if back == addrs[i] {
			return nil
		}
	}

	c := NewClient(back)

	prefix := colon.Escape(bin_name) + "::"

	tmp := trib.List{}
	e := c.ListKeys(&(trib.Pattern{prefix, ""}), &tmp)

	for _, key := range tmp.L {
		//Get list
		values := trib.List{}
		e = c.ListGet(key, &values)
		if e != nil {
			return e
		}

		for _, value := range values.L {
			n := 0
			c.ListRemove(&trib.KeyValue{key, value}, &n)
		}
	}
	return nil
}
