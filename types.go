package main

import (
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type factory interface {
	make(line string) task
}

type task interface {
	do() // do the task

	// output the result of the finished task to stdout / stderr.
	output()
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
//
// output() is defined in output.go and windows.go
type tryTask struct {
	addr    string        // address of remote host
	pass    string        // the password to try
	user    string        // user to use
	timeout time.Duration // how long to wait for a response

	result string // status of the password try
}

func (t *tryTask) do() {
	config := &ssh.ClientConfig{
		User:            t.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(t.pass)},
		Timeout:         t.timeout,
	}

retry:
	// Try and connect with the password used.
	_, err := ssh.Dial("tcp", t.addr, config)
	if err == nil {
		t.result = "ACCESS GRANTED"
		return
	}

	t.result = "FAILED"
	switch e := err.(type) {
	case net.Error:
		if e.Timeout() {
			t.result = "TIMEOUT"
		}

		log.Printf("net.Error: %v", e)
	default:
		if !strings.Contains(err.Error(), "unable to authenticate") {
			log.Printf("error: %q; retrying...", err)
			goto retry
		}
		log.Printf("%v", err)
	}
}
