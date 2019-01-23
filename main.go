// TODO
// Instead of options for username etc have it in args in the form user@adress:port
// Username list?
// If ssh password auth is not supported, detect it and stop the program

// Package main implements a simple ssh bruteforce tool to use with wordlists.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/UlisseMini/clean"
)

var (
	wordlist = flag.String("w", "wordlist.txt", "indicate wordlist file to use")
	addr     = flag.String("a", "127.0.0.1:22",
		"indicate the target address")

	user = flag.String("u", "root", "indicate user to use")
	// Set the timeout depending on the latency between you and the remote host.
	timeout = flag.Duration("t", 300*time.Millisecond,
		"set timeout for ssh dial response. do not set this too low!")

	debug = flag.Bool("d", false, "debug mode, print logs to stderr")
)

func main() {
	defer clean.Do()

	flag.Parse()
	printUsedValues()
	if *debug {
		log.SetOutput(os.Stderr)
	}

	passFile, err := os.Open(*wordlist)
	if err != nil {
		fmt.Printf("Error opening wordlist: %v\n", err)
		return
	}

	// add closing the file to global cleanup
	clean.Add(func() {
		passFile.Close()
	}, "passFile")

	scanner := bufio.NewScanner(passFile)

	// create factory
	fac := &sshFactory{
		user:    *user,
		timeout: *timeout,
		addr:    *addr,
	}

	// Get finished tasks from the finished channel
	finished := make(chan task)
	go func() {
		log.Println("starting to recv from finished")
		for t := range finished {
			// t.output will terminate the progam if it gets the right password
			t.output()
		}
		log.Println("done recv from finished")
	}()

	var wg sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("create task for %q", line)
		t := fac.make(line)

		wg.Add(1)
		go func(t task) {
			defer wg.Done()
			t.do()
			finished <- t
		}(t)
	}

	// wait for the tasks to be done
	log.Println("waiting for WaitGroup")
	wg.Wait()

	log.Println("closing finished")
	close(finished) // all tasks are finished so close finished.
}

func printUsedValues() {
	fmt.Printf("target: %s@%s\n", *user, *addr)
	fmt.Printf("timeout: %v\n", timeout)
	fmt.Printf("wordlist: %s\n", *wordlist)
}
