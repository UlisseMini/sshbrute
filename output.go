// +build !windows

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/UlisseMini/clean"
	"github.com/fatih/color"
)

func (t *tryTask) output() {
	pass := color.BlueString(t.pass)
	out := fmt.Sprintf("%s ", pass)

	switch t.result {
	case "ACCESS GRANTED":
		out += color.GreenString(t.result)
		fmt.Fprintf(os.Stdout, "%s\n", out)
		clean.Exit(0)
	case "FAILED":
		out += color.RedString(t.result)
	case "TIMEOUT":
		out += color.YellowString(t.result)
	default: // should never happen
		log.Printf("unknown result: %q", t.result)
		out += t.result
	}

	fmt.Fprintf(os.Stderr, "%s\n", out)
}
