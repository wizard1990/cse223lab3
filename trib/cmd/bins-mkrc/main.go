package main

import (
	"flag"
	"fmt"
	"log"

	"trib"
	"trib/randaddr"
)

var (
	local = flag.Bool("local", false, "always use local ports")
	nback = flag.Int("nback", 1, "number of back-ends")
	nkeep = flag.Int("nkeep", 1, "number of keepers")
	frc   = flag.String("rc", trib.DefaultRCPath, "bin storage config file")
	full  = flag.Bool("full", false, "setup of 10 back-ends and 3 keepers")
)

func main() {
	flag.Parse()

	if *nback > 10 {
		log.Fatal(fmt.Errorf("too many back-ends"))
	}
	if *nkeep > 10 {
		log.Fatal(fmt.Errorf("too many keepers"))
	}

	if *full {
		*nback = 10
		*nkeep = 3
	}

	p := randaddr.RandPort()

	rc := new(trib.RC)
	rc.Backs = make([]string, *nback)
	rc.Keepers = make([]string, *nkeep)

	if !*local {
		const ipOffset = 211
		for i := 0; i < *nback; i++ {
			host := fmt.Sprintf("172.22.14.%d", ipOffset+i)
			rc.Backs[i] = fmt.Sprintf("%s:%d", host, p)
		}

		for i := 0; i < *nkeep; i++ {
			host := fmt.Sprintf("172.22.14.%d", ipOffset+i)
			rc.Keepers[i] = fmt.Sprintf("%s:%d", host, p+1)
		}
	} else {
		for i := 0; i < *nback; i++ {
			rc.Backs[i] = fmt.Sprintf("localhost:%d", p)
			p++
		}

		for i := 0; i < *nkeep; i++ {
			rc.Keepers[i] = fmt.Sprintf("localhost:%d", p)
			p++
		}
	}

	fmt.Println(rc.String())

	if *frc != "" {
		e := rc.Save(*frc)
		if e != nil {
			log.Fatal(e)
		}
	}
}
