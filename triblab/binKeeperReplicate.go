package triblab

import(
	"fmt"
	"trib"
)

func (self *binKeeper) Replicate_bin() error{
	index := 0
	suffix := "::user"

  for{
		backend := self.clientMap[self.backs[index]]
		users := trib.List{[]string{}}

    e := backend.Keys(&(trib.Pattern{Suffix:suffix}),&users)
		if e == nil{
		  e = self.updateAll(users.L)
		  if e != nil{
		  }
		}

		index ++
		if index >= len(self.backs){
			index = 0
		}
	}
	return fmt.Errorf("replication stops for strange reasons")
}



func (self *binKeeper) update(key string, to trib.Storage, from trib.Storage) error{

//  toRemove := trib.List{[]string{}}
	source := trib.List{[]string{}}
	succ := false
	/*
	e := to.ListGet(key,&toRemove)
	if e != nil{
		return e
	}
	
	//remove everything in the invalid bin
	for _,toRemove := range toRemove.L{
		n := 0
    e = to.ListRemove(&trib.KeyValue{key,toRemove}, &n)
		if e != nil{
			return e
		}
	}
	*/

	//move everything from valid bin to invalid bin
	e := from.ListGet(key,&source)
	if e != nil{
		return e
	}
	for _,toMove := range source.L{
    e = to.ListAppend(&trib.KeyValue{key,toMove},&succ)
		if e != nil{
			return e
		}
	}
	return nil
}


func (self *binKeeper) updateAll(userKey []string) error{

  for _,binName := range userKey{
    binName = binName[:len(binName) - 7]
		if self.start_audit_bin(binName) == false{
			self.end_audit_bin(binName)
			continue
		}
		binsToAudit := self.bc.Bin(binName)
//binsToAudit := []trib.Storage{}
		validBins := []trib.Storage{}
		invalidBins := []trib.Storage{}

		for _,b := range binsToAudit{
			var v string
			e := b.Get("completed",&v)
			if e != nil{
				continue
			}
			if v == "1"{
				validBins = append(validBins,b)
			} else{
				invalidBins = append(invalidBins,b)
			}
		}

		if len(validBins) == 0{
	    fmt.Println("No valid bins, bin name :" + binName)
			continue
		}

		//replicate from valid bins to invalid bins
		for _,receiver := range invalidBins{
			// update following
			e := self.update("Following", receiver, validBins[0])
			if e != nil{
				continue
			}
			//update tribs
			e = self.update("Post", receiver, validBins[0])
			if e != nil{
				continue
			}
			//update username key-value
      succ := false
			receiver.Set(&trib.KeyValue{"user","1"}, &succ)

			if self.end_audit_bin(binName) == false{
				continue
			}
			receiver.Set(&trib.KeyValue{"completed","1"}, &succ)
		}
	}
	return nil

}
