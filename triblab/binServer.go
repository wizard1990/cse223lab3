package triblab

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
	"trib"
)

type binServer struct {
	server    trib.BinStorage
	lock      sync.Mutex
	userCache []string
}

//search user from bins
func (self *binServer) findUser(user string) (string, error) {
	for _, cachedUser := range self.userCache {
		if user == cachedUser {
			return "1", nil
		}
	}
	clients := self.server.Bin(user)
	res := ""
	var res string
	index := 0
	clients[0].Get("Completed", &res)
	if len(res) != 1 {
		clients[1].Get("Completed", &res)
		index = 1
		if len(res) != 1 {
			clients[2].Get("Completed", &res)
			index = 2
		}
	}
	if e := clients[index].Get(user, &res); e != nil {
		return "", e
	}
	if len(res) == 0 {
		return "", fmt.Errorf("username %q not exists", user)
	}
	return res, nil
}

//search user from user-list, used for sign up to keep consistency
func (self *binServer) findUserFromList(user string) (bool, error) {
	userList, e := self.ListAllUsers()
	if e != nil {
		return false, e
	}
	for _, username := range userList {
		if user == username {
			return true, nil
		}
	}
	return false, nil
}

func (self *binServer) SignUp(user string) error {
	if !trib.IsValidUsername(user) {
		return fmt.Errorf("invalid username %q", user)
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	found, e := self.findUserFromList(user)
	if e != nil {
		return e
	}
	if found {
		return fmt.Errorf("user %q already exists", user)
	}

	clients := self.server.Bin(user)
	succ := false
	if e := clients[0].Set(&trib.KeyValue{user, "1"}, &succ); e != nil {
		return e
	}
	if e := clients[1].Set(&trib.KeyValue{user, "1"}, &succ); e != nil {
		return e
	}
	if e := clients[2].Set(&trib.KeyValue{user, "1"}, &succ); e != nil {
		return e
	}
	if e := clients[0].Set(&trib.KeyValue{"Completed", "1"}, &succ); e != nil {
		return e
	}
	if e := clients[1].Set(&trib.KeyValue{"Completed", "1"}, &succ); e != nil {
		return e
	}
	if e := clients[2].Set(&trib.KeyValue{"Completed", "1"}, &succ); e != nil {
		return e
	}

	clients = self.server.Bin("ListUsers")
	succ = false
	e = clients[0].ListAppend(&trib.KeyValue{"ListUsers", user}, &succ)
	e = clients[1].ListAppend(&trib.KeyValue{"ListUsers", user}, &succ)
	e = clients[2].ListAppend(&trib.KeyValue{"ListUsers", user}, &succ)
	if (e == nil) && (len(self.userCache) < 20) {
		self.userCache = append(self.userCache, user)
	}
	return e
}

func (self *binServer) ListAllUsers() ([]string, error) {
	client := self.server.Bin("ListUsers")
	userList := trib.List{L: []string{}}
	e := client[0].ListGet("ListUsers", &userList)
	return userList.L, e
}

func (self *binServer) ListUsers() ([]string, error) {
	if len(self.userCache) == 20 {
		return self.userCache, nil
	}
	userList, e := self.ListAllUsers()
	if e != nil {
		return userList, e
	}
	if len(userList) > 20 {
		userList = userList[:20]
	}
	if len(userList) > len(self.userCache) {
		self.userCache = make([]string, len(userList))
		copy(self.userCache, userList)
	}
	return userList, nil
}

func (self *binServer) Post(who, post string, clock uint64) error {
	if len(post) > trib.MaxTribLen {
		return fmt.Errorf("trib too long")
	}
	if _, e := self.findUser(who); e != nil {
		return e
	}

	clients := self.server.Bin(who)
	var c uint64
	clients[0].Clock(clock, &c)
	clients[1].Clock(clock, &c)
	clients[2].Clock(clock, &c)

	newTrib := trib.Trib{
		User:    who,
		Message: post,
		Time:    time.Now(),
		Clock:   c,
	}
	b, e := json.Marshal(newTrib)
	if e != nil {
		return e
	}
	v := string(b)
	succ := false
	clients[0].ListAppend(&trib.KeyValue{"Post", v}, &succ)
	clients[1].ListAppend(&trib.KeyValue{"Post", v}, &succ)
	clients[2].ListAppend(&trib.KeyValue{"Post", v}, &succ)
	return nil
}

func (self *binServer) Tribs(user string) ([]*trib.Trib, error) {
	if _, e := self.findUser(user); e != nil {
		return nil, e
	}
	clients := self.server.Bin(user)
	plist := trib.List{L: []string{}}
	var res string
	index := 0
	clients[0].Get("Completed", &res)
	if len(res) != 1 {
		clients[1].Get("Completed", &res)
		index = 1
		if len(res) != 1 {
			clients[2].Get("Completed", &res)
			index = 2
		}
	}
	if e := clients[index].ListGet("Post", &plist); e != nil {
		return nil, e
	}

	tribList := make([]*trib.Trib, len(plist.L))
	dec := json.NewDecoder(strings.NewReader(strings.Join(plist.L, "")))
	for i := 0; ; i++ {
		var m trib.Trib
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		tribList[i] = &m
	}
	tribSort(tribList)
	l := len(tribList)
	if l > trib.MaxTribFetch {
		for i := trib.MaxTribFetch; i < l; i++ {
			go func(bc trib.Storage, old *trib.Trib) {
				b, e := json.Marshal(*old)
				if e != nil {
					return
				}
				v := string(b)
				n := 0
				bc.ListRemove(&trib.KeyValue{"Post", v}, &n)
			}(clients[index], tribList[i])
		}
	}
	return tribList[:l], nil
}

func (self *binServer) Follow(who, whom string) error {
	if who == whom {
		return fmt.Errorf("You cannot follow yourself, narcissist.")
	}
	if _, e := self.findUser(who); e != nil {
		return e
	}
	if _, e := self.findUser(whom); e != nil {
		return e
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	b, e := self.IsFollowing(who, whom)
	if e != nil {
		return e
	}
	if b {
		return fmt.Errorf("%q has already followed %q", who, whom)
	}

	clients := self.server.Bin(who)
	succ := false
	clients[0].ListAppend(&trib.KeyValue{"Following", whom}, &succ)
	clients[1].ListAppend(&trib.KeyValue{"Following", whom}, &succ)
	clients[2].ListAppend(&trib.KeyValue{"Following", whom}, &succ)
	return nil
}

func (self *binServer) Unfollow(who, whom string) error {
	if _, e := self.findUser(who); e != nil {
		return e
	}
	if _, e := self.findUser(whom); e != nil {
		return e
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	b, e := self.IsFollowing(who, whom)
	if e != nil {
		return e
	}
	if !b {
		return fmt.Errorf("%q has not followed %q yet", who, whom)
	}

	clients := self.server.Bin(who)
	n := 0
	clients[0].ListRemove(&trib.KeyValue{"Following", whom}, &n)
	clients[1].ListRemove(&trib.KeyValue{"Following", whom}, &n)
	clients[2].ListRemove(&trib.KeyValue{"Following", whom}, &n)
	return nil
}

func (self *binServer) IsFollowing(who, whom string) (bool, error) {
	if who == whom {
		return false, nil
	}
	if _, e := self.findUser(who); e != nil {
		return false, e
	}
	if _, e := self.findUser(whom); e != nil {
		return false, e
	}

	fs, e := self.Following(who)
	if e != nil {
		return false, e
	}
	for _, f := range fs {
		if f == whom {
			return true, nil
		}
	}
	return false, nil
}

func (self *binServer) Following(who string) ([]string, error) {
	if _, e := self.findUser(who); e != nil {
		return []string{}, e
	}
	clients := self.server.Bin(who)
	var res string
	index := 0
	clients[0].Get("Completed", &res)
	if len(res) != 1 {
		clients[1].Get("Completed", &res)
		index = 1
		if len(res) != 1 {
			clients[2].Get("Completed", &res)
			index = 2
		}
	}
	fs := trib.List{L: []string{}}
	e := clients[index].ListGet("Following", &fs)
	return fs.L, e
}

func (self *binServer) Home(user string) ([]*trib.Trib, error) {
	tribList := []*trib.Trib{}
	followList, e := self.Following(user)
	if e != nil {
		return tribList, e
	}
	followList = append(followList, user)
	//fmt.Println(len(followList))
	tribCh := make(chan []*trib.Trib, len(followList))
	for _, userName := range followList {
		go func(user string) {
			tbList, e := self.Tribs(user)
			//fmt.Println("tribs error", e)
			if e != nil {
				tribCh <- nil
			} else {
				tribCh <- tbList
			}
		}(userName)
	}
	maxHeap := tribHeap{&ByTime{}, trib.MaxTribFetch}
	//fmt.Println("start to push into heap")
	for i := 0; i < len(followList); i++ {
		tList := <-tribCh
		//fmt.Println(tList)
		if tList != nil {
			for _, tb := range tList {
				//fmt.Println("push!")
				heap.Push(&maxHeap, tb)
			}
		}
	}
	//fmt.Println("start to pop from heap")
	for len(*maxHeap.sorter) > 0 {
		newTrib := heap.Pop(&maxHeap)
		//fmt.Println(newTrib)
		tribList = append(tribList, newTrib.(*trib.Trib))
	}
	//fmt.Println("start clock sync")
	if len(tribList) > 0 {
		tribSort(tribList)
		newestTrib := tribList[len(tribList)-1]
		var n uint64 = 0
		clients := self.server.Bin(user)
		clients[0].Clock(newestTrib.Clock, &n)
		clients[1].Clock(newestTrib.Clock, &n)
		clients[2].Clock(newestTrib.Clock, &n)
	}
	return tribList, nil
}
