package tribtest

import (
	"runtime/debug"
	"sort"
	"testing"

	"trib"
)

func CheckServer(t *testing.T, server trib.Server) {
	ne := func(e error) {
		if e != nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	er := func(e error) {
		if e == nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
			debug.PrintStack()
			t.Fatal()
		}
	}

	ne(server.SignUp("h8liu"))
	er(server.SignUp(" h8liu"))
	er(server.SignUp("8hliu"))
	er(server.SignUp("H8liu"))

	ne(server.SignUp("fenglu"))

	users, e := server.ListUsers()
	ne(e)

	as(len(users) == 2)
	sort.Strings(users)
	as(users[0] == "fenglu")
	as(users[1] == "h8liu")

	ne(server.Follow("h8liu", "fenglu"))
	b, e := server.IsFollowing("h8liu", "fenglu")
	ne(e)
	as(b)

	b, e = server.IsFollowing("fenglu", "h8liu")
	ne(e)
	as(!b)

	b, e = server.IsFollowing("h8liu", "fenglu2")
	er(e)
	as(!b)

	ne(server.Unfollow("h8liu", "fenglu"))
	er(server.Unfollow("h8liu", "fenglu"))

	b, e = server.IsFollowing("h8liu", "fenglu")
	ne(e)
	as(!b)

	ne(server.Follow("h8liu", "fenglu"))

	clk := uint64(0)

	er(server.Post("", "", clk))

	longMsg := ""
	for i := 0; i < 200; i++ {
		longMsg += " "
	}

	er(server.Post("h8liu", longMsg, clk))
	ne(server.Post("h8liu", "hello, world", clk))

	clk = uint64(0)

	tribs, e := server.Tribs("h8liu")
	ne(e)
	as(len(tribs) == 1)
	tr := tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")
	if tr.Clock > clk {
		clk = tr.Clock
	}

	tribs, e = server.Home("fenglu")
	ne(e)
	as(tribs != nil)
	as(len(tribs) == 0)

	ne(server.Follow("fenglu", "h8liu"))
	tribs, e = server.Home("fenglu")
	ne(e)
	as(len(tribs) == 1)
	tr = tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")
	if tr.Clock > clk {
		clk = tr.Clock
	}

	ne(server.Post("h8liu", "hello, world2", clk))
	tribs, e = server.Home("fenglu")
	ne(e)
	as(len(tribs) == 2)
	tr = tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")

	tr = tribs[1]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world2")

	er(server.Follow("fenglu", "fenglu"))
	er(server.Follow("fengl", "fenglu"))
	er(server.Follow("fenglu", "fengl"))
	er(server.Follow("fenglu", "h8liu"))

	tribs, e = server.Home("h8liu")
	ne(e)
	as(len(tribs) == 2)
	tr = tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")

	tr = tribs[1]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world2")

	ne(server.SignUp("rkapoor"))
	fos, e := server.Following("rkapoor")
	ne(e)
	as(fos != nil)
	as(len(fos) == 0)
}
