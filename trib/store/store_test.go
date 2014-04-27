package store_test

import (
	"testing"

	. "trib/store"
	"trib/tribtest"
)

func TestStorage(t *testing.T) {
	s := NewStorage()

	tribtest.CheckStorage(t, s)
}
