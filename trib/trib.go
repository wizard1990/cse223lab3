// Package trib defines basic interfaces and constants
// for Tribbler service implementation.
package trib

import (
	"time"
)

const (
	MaxUsernameLen = 15   // Maximum length of a username
	MaxTribLen     = 140  // Maximum length of a tribble
	MaxTribFetch   = 100  // Maximum count of tribbles for Home() and Tribs()
	MinListUser    = 20   // Minimum count of users required for ListUsers()
	MaxFollowing   = 2000 // Maximum count of users that one can follow
)

type Trib struct {
	User    string    // who posted this trib
	Message string    // the content of the trib
	Time    time.Time // the physical timestamp
	Clock   uint64    // the logical timestamp
}

type Server interface {
	// Creates a user
	SignUp(user string) error

	// List 20 registered users.  When there are less than 20 users that
	// signed up the service, all of them needs to be listed.  When there
	// are more than 20 users that signed up the service, an arbitrary set
	// of at lest 20 of them needs to be listed.
	ListUsers() ([]string, error)

	// Post a tribble.  The clock is the maximum clock value this user has
	// seen so far by reading tribbles or clock sync.
	Post(who, post string, clock uint64) error

	// List the tribs that a particular user posted
	// The result should be sorted in alphabetical order
	Tribs(user string) ([]*Trib, error)

	// Follow someone's timeline
	Follow(who, whom string) error

	// Unfollow
	Unfollow(who, whom string) error

	// Returns true when who following whom
	IsFollowing(who, whom string) (bool, error)

	// Returns the list of following users
	Following(who string) ([]string, error)

	// List the trib of someone's following users
	Home(user string) ([]*Trib, error)
}

// Checks if a username is a valid one. Returns true if it is.
func IsValidUsername(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > MaxUsernameLen {
		return false
	}

	for i, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}

		if i > 0 && r >= '0' && r <= '9' {
			continue
		}

		return false
	}

	return true
}
