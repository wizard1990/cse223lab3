package triblab

import (
	"fmt"
	"strings"
	"time"
	"trib"
)

func (self *binKeeper) Replicate_bin() error {
	index := 0

	for {
		time.Sleep(time.Second * 1)
		backend := self.clientMap[self.backs[index]]
		users := trib.List{[]string{}}

		e := backend.ListKeys(&(trib.Pattern{Suffix: "KV"}), &users)
		if e == nil {
			self.updateAll(users.L, "KV")
		} else {
			fmt.Println(e)
		}

		/*
			e = backend.Keys(&(trib.Pattern{Suffix: "L"}), &users)
			if e == nil{
			  self.updateAll(users.L, "KV")
			}
		*/

		index++
		if index >= len(self.backs) {
			index = 0
		}
	}
	return fmt.Errorf("replication stops for strange reasons")
}

func (self *binKeeper) update(key string, bins []trib.Storage) error {

	lists := make([]trib.List, 3)
	for i, _ := range bins {
		bins[i].ListGet(key, &lists[i])
	}
	fmt.Println(lists[0], lists[1], lists[2])
	_, maxSet, _ := FindLargestClock(&lists[0], &lists[1], &lists[2])

	for i, origin := range lists {
		fmt.Println("from->to", maxSet, origin)
		toAdd := DiffList(maxSet, &origin)
		fmt.Println(toAdd)
		for _, listToAdd := range toAdd.L {
			succ := false
			bins[i].ListAppend(&trib.KeyValue{key, listToAdd}, &succ)
		}

	}
	return nil
}

func (self *binKeeper) updateAll(users []string, suffix string) error {

	for _, binName := range users {
		fmt.Println(binName)

		binName = strings.TrimRight(binName, "::"+suffix)
		fmt.Println("binName", binName) //Bug to be fixed..

		//Colon.Unescape ...and etc

		//TrimRight
		if self.start_audit_bin(binName) == false {
			self.end_audit_bin(binName)
			continue
		}
		binsToAudit := self.bc.KeeperBin(binName)
		//self.update(suffix, binsToAudit)
		//Temp for debuging
		self.update(binName+"::"+suffix, binsToAudit)

		//update username key-value
		self.end_audit_bin(binName)

	}
	return nil
}
