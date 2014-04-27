package triblab

import (
    "trib"
    "fmt"
    "sync"
    "encoding/json"
    "strings"
)

type binServer struct {
    server trib.BinStorage
    lock sync.Mutex
    userCache []string
}

//search user from bins
func (self *Server) findUser(user string) (state, error) {
    for _, cachedUser := range self.userCache {
        if user == cachedUser {
            return "1", nil
        }
    }
    client := self.server.Bin(user)
    res := ""
    if e := client.Get(user, &res); e != nil {
        return nil, e
    }
    if len(res) == 0 {
        return nil, fmt.Errorf("username %q not exists", user)
    }
    return res, nil
}

//search user from user-list, used for sign up to keep consistency
func (self *Server) findUserFromList(user string) (bool, error) {
    userList, e = self.ListAllUsers()
    if(e != nil){
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

    client := self.server.Bin(user)
    if e := client.Set(user, "1"); e != nil {
        return e
    }
    client = self.server.Bin("ListUsers")
    succ := false
    e = client.ListAppend(&trib.KeyValue{"ListUsers", user}, &succ)
    if (e == nil) && (len(self.userCache) < 20) {
        self.userCache = append(self.userCache, user)
    }
    return e
}

func (self *binServer) ListAllUsers() ([]string, error){
    client := self.server.Bin("ListUsers")
    userList := trib.List{L:[]string{}}
    e := client.ListGet("ListUsers", &userList)
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
        self.userCache = make([]string, len(userList)
        copy(self.userCache, userList)
    }
    return userList, nil
}

func (self *binServer) Post(who, post string, clock uint64) error {
    if len(post) > trib.MaxTribLen {
        return fmt.Errorf("trib too long")
    }
    if _, e = self.findUser(who); e != nil {
        return e
    }

    self.lock.Lock()
    defer self.lock.Unlock()

    client := self.server.Bin(who)
    var c uint64
    client.Clock(clock, &c)

    newTrib := trib.Trib {
        User: who, 
        Message: post, 
        Time: time.Now(), 
        Clock: c,
    }
    b, e := json.Marshal(newTrib)
    if e != nil {
        return e
    }
    v := string(b)
    succ := false
    return client.ListAppend(&trib.KeyValue{"Post", value}, &succ)
}

func (self *binServer) Tribs(user string) ([]*trib.Trib, error) {
    if _, e = self.findUser(who); e != nil {
        return nil, e
    }
    client := self.server.Bin(user)
    plist := trib.List{L:[]string{}}
    if e := client.ListGet("Post", &plist); e != nil {
        return nil, e
    }

    tribList = make([]*trib.Trib, len(plist.L))
    dec := json.NewDecoder(strings.NewReader(strings.Join(plist.L, "")))
    for i := 0;;i++{
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
            client.List
        }

    }
    
    return tribList[:l], nil
}

func (self *binServer) Follow(who, whom string) error {
    if(who == whom) {
        return fmt.Errorf("You cannot follow yourself, narcissist.")
    }
    if _, e = self.findUser(who); e != nil {
        return e
    }
    if _, e = self.findUser(whom); e != nil {
        return e
    }

    self.lock.Lock()
    defer self.lock.Unlock()
    
    b, e := self.IsFollowing(who, whom)
    if(e != nil){
        return e
    }
    if(flag){
        return fmt.Errorf("%q has already followed %q", who, whom)
    }
    
    client := self.server.Bin(who)
    succ := false
    return client.ListAppend(&trib.KeyValue{"Following", whom}, &succ)
}

func (self *binServer) Unfollow(who, whom string) error {
    if _, e = self.findUser(who); e != nil {
        return e
    }
    if _, e = self.findUser(whom); e != nil {
        return e
    }

    self.lock.Lock()
    defer self.lock.Unlock()
    
    b, e := self.IsFollowing(who, whom)
    if(e != nil){
        return e
    }
    if(!flag){
        return fmt.Errorf("%q has not followed %q yet", who, whom)
    }
    
    client := self.server.Bin(who)
    n := 0
    return client.ListRemove(&trib.KeyValue{"Following", whom}, &n)
}

func (self *binServer) IsFollowing(who, whom string) (bool, error) {
    if who == whom {
        return false, nil
    }
    if _, e = self.findUser(who); e != nil {
        return false, e
    }
    if _, e = self.findUser(whom); e != nil {
        return false, e
    }

    fs, e := self.Following(who)
    if(e != nil){
        return false, e
    }
    for _, f := range fs {
        if(f == whom) {
            return true, nil
        }
    }
    return false, nil
}

func (self *binServer) Following(who string) ([]string, error) {
    if _, e = self.findUser(who); e != nil {
        return false, e
    }
    client := self.server.Bin(who)
    fs := trib.List{L:[]string{}}
    e := client.ListGet("Following", &fs)
    return fs.L, e
}

func (self *binServer) Home(user string) ([]*Trib, error) {

}