# Distributed, Concurrent Issue Board

## Prerequisites

This project uses Logrus as a logging framework. To install Logrus, type this in your shell: `go get "github.com/sirupsen/logrus"`.

## Running

~~~bash
# on Unix-based systems:
go build && ./NoticeBoard [args]

# alternatively, you can use...
# this does not work on Windows
go run *.go [args]

# on Windows:
go build
.\NoticeBoard.exe [args]

# example:
go run *.go localhost:4007 localhost:4008 localhost:4009
~~~

Where `args` is a list of `IP:Port` pairs, with the first pair being the host (you).

## Using the system

When you first run the program you are greeted by a prompt asking you for your username.

