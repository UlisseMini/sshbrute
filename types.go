package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

type factory interface {
	make(line string) task
}

type task interface {
	do() // do the task

	// output the result of the finished task to stdout / stderr,
	// it returns true if all other tasks should stop.
	output() bool
}

type sshFactory struct {
	user    string        // username to use with ssh tries
	timeout time.Duration // timeout to pass to ssh.Dial
	addr    string
}

func (s *sshFactory) make(pass string) task {
	return &tryTask{
		addr:    s.addr,
		pass:    pass,
		result:  "UNFINISHED",
		user:    s.user,
		timeout: s.timeout,
	}
}

// tryTask implements the task interface,
// it tries a password with a username on the remote host.
type tryTask struct {
	addr    string        // address of remote host
	pass    string        // the password to try
	user    string        // user to use
	timeout time.Duration // how long to wait for a response

	result string // status of the password try
}

func (t *tryTask) output() (allDone bool) {
	pass := color.BlueString(t.pass)

	if t.result == "ACCESS GRANTED" {
		fmt.Fprintf(os.Stdout, "%s %s\n", pass, color.GreenString(t.result))
		allDone = true
	}

	if t.result == "FAILED" {
		fmt.Fprintf(os.Stderr, "%s %s\n", pass, color.RedString(t.result))
	}

	// Should never happen
	if t.result == "UNFINISHED" {
		fmt.Fprintf(os.Stderr, "%s %s\n", pass, color.YellowString(t.result))
	}

	return allDone
}

func (t *tryTask) do() {
	config := &ssh.ClientConfig{
		User:            t.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(t.pass)},
		Timeout:         t.timeout,
	}

	// Try and connect with the password used.
	_, err := ssh.Dial("tcp", t.addr, config)
	if err == nil {
		t.result = "ACCESS GRANTED"
		return
	}

	t.result = "FAILED"
}
