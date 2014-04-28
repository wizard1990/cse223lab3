package triblab

import(
	"fmt"
	"trib"
	"strings"
)

func (self *binKeeper) Replicate_bin() error{
	index := 0

  for{
		backend := self.clientMap[self.backs[index]]
		users := trib.List{[]string{}}

		e := backend.Keys(&(trib.Pattern{Suffix: "kv"}), &users)
		if e == nil {
			self.updateAll(users.L, "kv")
		}

		e = backend.Keys(&(trib.Pattern{Suffix: "L"}), &users)
		if e == nil{
		  self.updateAll(users.L, "kv")
		}

		index ++
		if index >= len(self.backs){
			index = 0
		}
	}
	return fmt.Errorf("replication stops for strange reasons")
}



func (self *binKeeper) update(key string, bins []trib.Storage) error{

	lists := make([]trib.List,3)
	for i,_ := range bins{
    bins[i].ListGet(key,&lists[i])
	}
	_,maxSet,_ := FindLargestClock(&lists[0],&lists[1],&lists[2])

	for i,origin := range lists{
		toAdd := DiffList(maxSet, &origin)
		for _,listToAdd := range toAdd.L{
			succ := false
      bins[i].ListAppend(&trib.KeyValue{key,listToAdd},&succ)
		}

	}
	return nil
}

func (self *binKeeper) updateAll(users []string, suffix string) error {

	for _, binName := range users {

		binName = strings.TrimRight(binName, "::" + suffix)
		//TrimRight
		if self.start_audit_bin(binName) == false {
			self.end_audit_bin(binName)
			continue
		}
		binsToAudit := self.bc.KeeperBin(binName)
		self.update(suffix, binsToAudit)

		//update username key-value
		self.end_audit_bin(binName)

	}
	return nil
}
