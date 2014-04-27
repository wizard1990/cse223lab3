package main

import (
	"log"
)

func errString(e error) string {
	if e == nil {
		return ""
	}
	return e.Error()
}

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func logError(e error) {
	if e != nil {
		log.Print(e)
	}
}
