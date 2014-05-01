package triblab

import (
	"fmt"
	"strings"
//	"time"
	"trib"
	"trib/colon"
)

func (self *binKeeper) Replicate_bin() error {
	index := 0

		for {
//		time.Sleep(time.Second * 1)
		backend := self.clientMap[self.backs[index]]
//fmt.Println(backend)
		users := trib.List{[]string{}}

e := backend.ListKeys(&(trib.Pattern{Suffix: "::1::KV"}), &users)

		if e == nil {
			 self.updateAll(users.L, backend)
		} else {
			fmt.Println(e)
		}

		index++
		if index >= len(self.backs) {
			index = 0
		}
	}
	return fmt.Errorf("replication stops for strange reasons")
}

func (self *binKeeper) update(key string, bins []trib.Storage) error {

	for _,b := range bins {
		pattern := trib.Pattern{Suffix: "::"+key}
		keysToCheck := trib.List{[]string{}}
		b.ListKeys(&pattern, &keysToCheck)
//fmt.Println("keys to check:", keysToCheck.L)
		for i,_ := range keysToCheck.L{
	    lists := make([]trib.List, 3)
			for j,_ := range bins{
				bins[j].ListGet(keysToCheck.L[i],&lists[j])
			}
			_, maxSet, _ := FindLargestClock(&lists[0], &lists[1], &lists[2])
//fmt.Println("maxset: ",maxSet)
			for j,origin := range lists{
        toAdd := DiffList(maxSet, &origin)
				for _, listToAdd := range toAdd.L{
          succ := false
//fmt.Println("listappend:",listToAdd)
					bins[j].ListAppend(&trib.KeyValue{keysToCheck.L[i], listToAdd}, &succ)
				}
			}
		}
	}
	return nil
}

func (self *binKeeper) updateAll(users []string, backend client) error {


	for _, binName := range users {


		binName = colon.Unescape(strings.TrimRight(binName, "::1::KV"))

		if self.start_audit_bin(binName) == false {
			self.end_audit_bin(binName)
			continue
		}
		binsToAudit, binsAddr := self.bc.KeeperBin(binName)
		self.update("KV", binsToAudit)
		self.update("L", binsToAudit)
//fmt.Println("about to deredundant", backend.addr, binsAddr, suffix)
		self.deRedundant(binName, backend, binsAddr, "KV")
		self.deRedundant(binName, backend, binsAddr, "L")

		//update username key-value
		self.end_audit_bin(binName)

	}
	return nil
}

func (self *binKeeper) deRedundant(binName string, backend client, binsAddr []string, suffix string) error{
  if backend.addr != binsAddr[0] && backend.addr != binsAddr[1] && backend.addr != binsAddr[2] {
    keys := trib.List{[]string{}}
		backend.ListKeys(&trib.Pattern{Prefix:binName, Suffix:suffix}, &keys)
    for _,key := range keys.L{
			values := trib.List{[]string{}}
			backend.ListGet(key, &values)
			for _,value := range values.L{
				n := 0
        backend.ListRemove(&trib.KeyValue{key, value}, &n)

//fmt.Println("about to remove ", key, value, n)
			}
		}
	}
	return nil
}
