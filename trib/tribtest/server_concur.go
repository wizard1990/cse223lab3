package tribtest

import (
	"runtime"
	"runtime/debug"
	"strconv"
	"testing"

	"trib"
)

func CheckServerConcur(t *testing.T, server trib.Server) {
	runtime.GOMAXPROCS(2)

	ne := func(e error) {
		if e != nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	er := func(e error) {
		if e == nil {
			debug.PrintStack()
			t.Fatal()
		}
	}

	as := func(cond bool) {
		if !cond {
			debug.PrintStack()
			t.Fatal()
		}
	}

	ne(server.SignUp("user"))

	p := func(th, n int, done chan<- bool) {
		for i := 0; i < n; i++ {
			ne(server.Post("user", strconv.Itoa(th*100+n), 0))
		}
		done <- true
	}

	nconcur := 9
	done := make(chan bool, nconcur)
	for i := 0; i < nconcur; i++ {
		go p(i, 10, done)
	}

	for i := 0; i < nconcur; i++ {
		<-done
	}

	ret, e := server.Tribs("user")
	ne(e)
	as(len(ret) == 10*nconcur || len(ret) == trib.MaxTribFetch)

	ne(server.SignUp("other"))
	fo := func(done chan<- bool) {
		e := server.Follow("user", "other")
		done <- (e == nil)
	}

	unfo := func(done chan<- bool) {
		e := server.Unfollow("user", "other")
		done <- (e == nil)
	}

	for i := 0; i < nconcur; i++ {
		go fo(done)
	}
	cnt := 0
	for i := 0; i < nconcur; i++ {
		if <-done {
			cnt++
		}
	}
	t.Logf("%d follows", cnt)
	// as(cnt == 1)

	er(server.Follow("user", "other"))

	fos, e := server.Following("user")
	ne(e)
	as(len(fos) == 1)
	as(fos[0] == "other")

	for i := 0; i < nconcur; i++ {
		go unfo(done)
	}
	cnt = 0
	for i := 0; i < nconcur; i++ {
		if <-done {
			cnt++
		}
	}
	t.Logf("%d unfollows", cnt)
	// as(cnt == 1)

	fos, e = server.Following("user")
	ne(e)
	as(len(fos) == 0)
}
