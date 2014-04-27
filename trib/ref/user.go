package ref

import (
	"sort"
	"time"
	"trib"
)

type user struct {
	following map[string]*user
	followers map[string]*user
	seqTribs  []*seqTrib
	tribs     []*trib.Trib
	home      []*trib.Trib
}

func newUser() *user {
	return &user{
		following: make(map[string]*user),
		followers: make(map[string]*user),
		seqTribs:  make([]*seqTrib, 0, 1024),
		tribs:     make([]*trib.Trib, 0, 1024),
		home:      make([]*trib.Trib, 0, 4096),
	}
}

func (self *user) isFollowing(whom string) bool {
	_, found := self.following[whom]
	return found
}

func (self *user) rebuildHome() {
	home := make([]*seqTrib, 0, 4096)
	home = append(home, self.seqTribs...)
	for _, user := range self.following {
		home = append(home, user.seqTribs...)
	}

	sort.Sort(bySeq(home))

	self.home = make([]*trib.Trib, 0, len(home))
	for _, t := range home {
		self.home = append(self.home, t.Trib)
	}
}

func (self *user) follow(whom string, u *user) {
	self.following[whom] = u
	self.rebuildHome()
}

func (self *user) unfollow(whom string) {
	delete(self.following, whom)
	self.rebuildHome()
}

func (self *user) addFollower(who string, u *user) {
	self.followers[who] = u
}

func (self *user) removeFollower(who string) {
	delete(self.followers, who)
}

func (self *user) listFollowing() []string {
	ret := make([]string, 0, len(self.following))
	for u := range self.following {
		ret = append(ret, u)
	}
	return ret
}

func (self *user) post(who, msg string, seq uint64, ts time.Time) {
	// make the new trib
	t := &trib.Trib{
		User:    who,
		Message: msg,
		Time:    ts,
		Clock:   seq,
	}

	// append a sequencial number, used in rebuilding subscribtion
	seqt := &seqTrib{
		seq:  seq,
		Trib: t,
	}

	// add to my own tribs
	self.tribs = append(self.tribs, t)
	self.seqTribs = append(self.seqTribs, seqt)

	// and it into the home timeline of my followers
	for _, user := range self.followers {
		user.home = append(user.home, t)
	}
	self.home = append(self.home, t)
}

func (self *user) listHome() []*trib.Trib {
	ntrib := len(self.home)
	start := 0
	if ntrib > trib.MaxTribFetch {
		start = ntrib - trib.MaxTribFetch
	}

	return self.home[start:]
}

func (self *user) listTribs() []*trib.Trib {
	ntrib := len(self.tribs)
	start := 0
	if ntrib > trib.MaxTribFetch {
		start = ntrib - trib.MaxTribFetch
	}
	return self.tribs[start:]
}
