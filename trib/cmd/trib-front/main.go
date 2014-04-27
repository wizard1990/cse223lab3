// Tribbler front end launcher program.
package main

import (
	"encoding/json"
	"flag"
	"go/build"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"trib"
	"trib/randaddr"
	"trib/ref"
	"triblab"
)

var (
	verbose = flag.Bool("v", false, "verbose logging")
	lab     = flag.Bool("lab", false, "use lab implementation")
	addr    = flag.String("addr", "localhost:rand", "serve address")
	frc     = flag.String("rc", trib.DefaultRCPath,
		"bin storage config file")
	dbinit = flag.Bool("init", false, "do not populate with test data")

	server trib.Server
)

func handleApi(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/")

	reply := func(obj interface{}) {
		bytes, e := json.Marshal(obj)
		noError(e)

		_, e = w.Write(bytes)
		logError(e)
	}

	bytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		return
	}
	input := string(bytes)

	if !*verbose {
		log.Println(name, input)
	}

	switch name {
	case "add-user":
		e = server.SignUp(input)
		if e != nil {
			reply(NewUserList(nil, e))
			break
		}
		ret, e := server.ListUsers()
		reply(NewUserList(ret, e))

	case "list-users":
		ret, e := server.ListUsers()
		reply(NewUserList(ret, e))

	case "list-tribs":
		tribs, e := server.Tribs(input)
		reply(NewTribList(tribs, e))

	case "list-home":
		tribs, e := server.Home(input)
		reply(NewTribList(tribs, e))

	case "is-following":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		ret, e := server.IsFollowing(ww.Who, ww.Whom)
		reply(NewBool(ret, e))

	case "follow":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Follow(ww.Who, ww.Whom)
		reply(NewBool(e == nil, e))

	case "unfollow":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Unfollow(ww.Who, ww.Whom)
		reply(NewBool(false, e))

	case "following":
		ret, e := server.Following(input)
		reply(NewUserList(ret, e))

	case "post":
		p := new(Post)
		e := json.Unmarshal(bytes, p)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Post(p.Who, p.Message, p.Clock)
		reply(NewBool(e == nil, e))

	default:
		w.WriteHeader(404)
	}
}

func ne(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func makeServer() trib.Server {
	if !*lab {
		return ref.NewServer()
	}

	rc, e := trib.LoadRC(*frc)
	ne(e)

	c := triblab.NewBinClient(rc.Backs)

	return triblab.NewFront(c)
}

func populate(server trib.Server) {
	ne(server.SignUp("h8liu"))
	ne(server.SignUp("fenglu"))
	ne(server.SignUp("rkapoor"))

	ne(server.Post("h8liu", "Hello, world.", 0))
	ne(server.Post("h8liu", "Just tribble it.", 0))
	ne(server.Post("fenglu", "Double tribble.", 0))
	ne(server.Post("rkapoor", "Triple tribble.", 0))

	ne(server.Follow("fenglu", "h8liu"))
	ne(server.Follow("fenglu", "rkapoor"))

	ne(server.Follow("rkapoor", "h8liu"))
}

func wwwPath() string {
	pkg, e := build.Import("trib", "./", build.FindOnly)
	if e != nil {
		log.Fatal(e)
	}
	return filepath.Join(pkg.Dir, "www")
}

func main() {
	flag.Parse()

	server = makeServer()
	if *dbinit {
		populate(server)
	}
	*addr = randaddr.Resolve(*addr)
	log.Printf("serve on %s", *addr)

	http.Handle("/", http.FileServer(http.Dir(wwwPath())))
	http.HandleFunc("/api/", handleApi)

	for {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
