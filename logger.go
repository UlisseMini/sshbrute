package main

import (
	"io/ioutil"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(ioutil.Discard)
}
