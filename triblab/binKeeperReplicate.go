package triblab

import (
	"fmt"
	"strings"
	"time"
	"trib"
	"trib/colon"
)

func (self *binKeeper) Replicate_bin() error {
	index := 0

	for {
		time.Sleep(time.Second * 1)
		backend := self.clientMap[self.backs[index]]
		users := trib.List{[]string{}}

//		e := backend.ListKeys(&(trib.Pattern{Suffix: "KV"}), &users)
e := backend.ListKeys(&(trib.Pattern{Suffix: "::1::KV"}), &users)

		if e == nil {
			self.updateAll(users.L, "KV")
			self.updateAll(users.L,"L")
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

	for _,b := range bins {
		pattern := trib.Pattern{Suffix: "::"+key}
		keysToCheck := trib.List{[]string{}}
		b.ListKeys(&pattern, &keysToCheck)
fmt.Println("keys to check:", keysToCheck.L)
		for i,_ := range keysToCheck.L{
	    lists := make([]trib.List, 3)
			for j,_ := range bins{
				bins[j].ListGet(keysToCheck.L[i],&lists[j])
			}
			_, maxSet, _ := FindLargestClock(&lists[0], &lists[1], &lists[2])
fmt.Println("maxset: ",maxSet)
			for j,origin := range lists{
        toAdd := DiffList(maxSet, &origin)
				for _, listToAdd := range toAdd.L{
          succ := false
fmt.Println("listappend:",listToAdd)
					bins[j].ListAppend(&trib.KeyValue{keysToCheck.L[i], listToAdd}, &succ)
				}
			}
		}
	}
	/*
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
	*/
	return nil
}

func (self *binKeeper) updateAll(users []string, suffix string) error {

//fmt.Println("updateALL")

	for _, binName := range users {
		time.Sleep(time.Microsecond * 100)
		fmt.Println(binName)

		binName = colon.Unescape(strings.TrimRight(binName, "::1::KV"))
//		fmt.Println("binName: ", binName) //Bug to be fixed..

		//Colon.Unescape ...and etc

		//TrimRight
		if self.start_audit_bin(binName) == false {
			self.end_audit_bin(binName)
			continue
		}
		binsToAudit := self.bc.KeeperBin(binName)
//fmt.Println("KeeperBin returns: ", binsToAudit)
		//self.update(suffix, binsToAudit)
		//Temp for debuging
		self.update(suffix, binsToAudit)

		//update username key-value
		self.end_audit_bin(binName)

	}
	return nil
}
