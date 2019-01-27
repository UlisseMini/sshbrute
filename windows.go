// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/UlisseMini/clean"
)

func (t *tryTask) output() {
	fmt.Fprintf(os.Stderr, "%s %s\n", t.pass, t.result)

	if t.result == "ACCESS GRANTED" {
		clean.Exit(0)
	}
}
