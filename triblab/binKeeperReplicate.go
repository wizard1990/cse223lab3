package triblab

import (
	"fmt"
	"trib"
)

func (self *binKeeper) Replicate_bin() error {
	index := 0
	suffix := "::kv"

	for {
		backend := self.clientMap[self.backs[index]]
		users := trib.List{[]string{}}

		e := backend.Keys(&(trib.Pattern{Suffix: suffix}), &users)
		if e == nil {
			e = self.updateAll(users.L)
		}

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
	_, maxSet, _ := FindLargestClock(&lists[0], &lists[1], &lists[2])

	for i, origin := range lists {
		toAdd := DiffList(maxSet, &origin)
		for _, listToAdd := range toAdd.L {
			succ := false
			bins[i].ListAppend(&trib.KeyValue{key, listToAdd}, &succ)
		}

	}
	return nil
}

func (self *binKeeper) updateAll(users []string) error {

	for _, binName := range users {
		binName = binName[:len(binName)-5]
		//TrimRight
		if self.start_audit_bin(binName) == false {
			self.end_audit_bin(binName)
			continue
		}
		binsToAudit := self.bc.KeeperBin(binName)
		self.update("kv", binsToAudit)
		self.update("L", binsToAudit)

		//update username key-value
		self.end_audit_bin(binName)

	}
	return nil
}
