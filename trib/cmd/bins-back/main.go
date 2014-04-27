// Tribbler back-end launcher.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"trib"
	"trib/local"
	"trib/store"
	"triblab"
)

var (
	frc = flag.String("rc", trib.DefaultRCPath, "bin storage config file")
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	flag.Parse()

	rc, e := trib.LoadRC(*frc)
	noError(e)

	run := func(i int) {
		if i > len(rc.Backs) {
			noError(fmt.Errorf("back-end index out of range: %d", i))
		}

		backConfig := rc.BackConfig(i, store.NewStorage())
		log.Printf("bin storage back-end serving on %s", backConfig.Addr)
		noError(triblab.ServeBack(backConfig))
	}

	args := flag.Args()

	n := 0
	if len(args) == 0 {
		// scan for addresses on this machine
		for i, b := range rc.Backs {
			if local.Check(b) {
				go run(i)
				n++
			}
		}

		if n == 0 {
			log.Fatal("no back-end found for this host")
		}
	} else {
		// scan for indices for the addresses
		for _, a := range args {
			i, e := strconv.Atoi(a)
			noError(e)
			go run(i)
			n++
		}
	}

	if n > 0 {
		select {}
	}
}
