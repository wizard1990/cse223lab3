package main

import (
	"trib"
)

type UserList struct {
	Err   string
	Users []string
}

type TribList struct {
	Err   string
	Tribs []*trib.Trib
}

type Bool struct {
	Err string
	V   bool
}

type Clock struct {
	Err string
	N   uint64
}

func NewTribList(tribs []*trib.Trib, e error) *TribList {
	return &TribList{errString(e), tribs}
}

func NewUserList(users []string, e error) *UserList {
	return &UserList{errString(e), users}
}

func NewBool(b bool, e error) *Bool {
	return &Bool{errString(e), b}
}

func NewClock(c uint64, e error) *Clock {
	return &Clock{errString(e), c}
}

type WhoWhom struct {
	Who  string
	Whom string
}

type Post struct {
	Who     string
	Message string
	Clock   uint64
}
