package ref_test

import (
	"testing"

	. "trib/ref"
	"trib/tribtest"
)

func TestTrib(t *testing.T) {
	server := NewServer()
	tribtest.CheckServer(t, server)

	server = NewServer()
	tribtest.CheckServerConcur(t, server)
}
