// Tribbler back-end keeper launcher.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"trib"
	"trib/local"
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
		if i > len(rc.Keepers) {
			noError(fmt.Errorf("keeper index out of range: %d", i))
		}

		keeperConfig := rc.KeeperConfig(i)
		c := make(chan bool)
		keeperConfig.Ready = c
		go func() {
			noError(triblab.ServeKeeper(keeperConfig))
		}()

		b := <-c
		if b {
			log.Printf("bin storage keeper serving on %s",
				keeperConfig.Addr())
		} else {
			log.Printf("bin storage keeper on %s init failed",
				keeperConfig.Addr())
		}
	}

	args := flag.Args()
	n := 0
	if len(args) == 0 {
		for i, k := range rc.Keepers {
			if local.Check(k) {
				go run(i)
				n++
			}
		}

		if n == 0 {
			log.Fatal("no keeper found for this host")
		}
	} else {
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
