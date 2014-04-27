package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"trib"
	"triblab"
)

func noError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func logError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
}

func kv(k, v string) *trib.KeyValue {
	return &trib.KeyValue{k, v}
}

func pat(pre, suf string) *trib.Pattern {
	return &trib.Pattern{pre, suf}
}

func kva(args []string) *trib.KeyValue {
	if len(args) == 1 {
		return kv("", "")
	} else if len(args) == 2 {
		return kv(args[1], "")
	}
	return kv(args[1], args[2])
}

func pata(args []string) *trib.Pattern {
	if len(args) == 1 {
		return pat("", "")
	} else if len(args) == 2 {
		return pat(args[1], "")
	}
	return pat(args[1], args[2])
}

func single(args []string) string {
	if len(args) == 1 {
		return ""
	}
	return args[1]
}

func printList(lst trib.List) {
	for _, e := range lst.L {
		fmt.Println(e)
	}
}

const cmdHelp = `Command List:
   bin [<bin name>]
   get <key>
   set <key> <value>
   keys [<prefix> [<suffix>]]
   list-get <key>
   list-append <key> <value>
   list-remove <key> <value>
   list-keys [<prefix> [<suffix]]
   clock [<atleast=0>]
   help
   exit
`

type client struct {
	bin string
	s   trib.BinStorage
}

func (self *client) printBin() {
	fmt.Printf("(working on bin %q)\n", self.bin)
}

func (self *client) runCmd(args []string) bool {
	var v string
	var b bool
	var lst trib.List
	var n int
	var cret uint64

	s := self.s.Bin(self.bin)

	cmd := args[0]

	switch cmd {
	case "get":
		logError(s.Get(single(args), &v))
		fmt.Println(v)
	case "set":
		logError(s.Set(kva(args), &b))
		fmt.Println(b)
	case "keys":
		logError(s.Keys(pata(args), &lst))
		printList(lst)
	case "list-get":
		logError(s.ListGet(single(args), &lst))
		printList(lst)
	case "list-append":
		logError(s.ListAppend(kva(args), &b))
		fmt.Println(b)
	case "list-remove":
		logError(s.ListRemove(kva(args), &n))
		fmt.Println(n)
	case "list-keys":
		logError(s.ListKeys(pata(args), &lst))
		printList(lst)
	case "clock":
		var c uint64
		var e error
		if len(args) >= 2 {
			c, e = strconv.ParseUint(args[1], 10, 64)
			logError(e)
		}
		logError(s.Clock(c, &cret))
		fmt.Println(cret)
	case "bin":
		if len(args) > 1 {
			self.bin = args[1]
		}
		self.printBin()
	case "help":
		fmt.Println(cmdHelp)
	case "exit":
		return true
	default:
		logError(fmt.Errorf("bad command, try \"help\"."))
	}
	return false
}

func fields(s string) []string {
	return strings.Fields(s)
}

func (self *client) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	self.printBin()

	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()
		args := fields(line)
		if len(args) > 0 {
			if self.runCmd(args) {
				break
			}
		}
		fmt.Print("> ")
	}

	e := scanner.Err()
	if e != nil {
		panic(e)
	}
}

var (
	frc = flag.String("rc", trib.DefaultRCPath, "bin storage config file")
)

func main() {
	flag.Parse()

	rc, e := trib.LoadRC(*frc)
	noError(e)

	s := triblab.NewBinClient(rc.Backs)

	cl := &client{s: s}
	cl.runPrompt()
	fmt.Println()
}
