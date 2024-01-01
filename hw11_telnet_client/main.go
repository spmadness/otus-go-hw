package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultDuration = "10s"

var ErrWrongDurationArg = errors.New("wrong duration string value")

var durationArg string

func main() {
	flag.StringVar(&durationArg, "timeout", defaultDuration, "file to read from")
	flag.Parse()

	duration, err := time.ParseDuration(durationArg)
	if err != nil {
		fmt.Printf("duration arg parse error: %s\n", ErrWrongDurationArg)
		os.Exit(1)
	}

	if len(flag.Args()) != 2 {
		fmt.Println("go-telnet usage: go-envdir [--timeout=10s] host port")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGPIPE)

	address := net.JoinHostPort(flag.Args()[0], flag.Args()[1])

	client := NewTelnetClient(address, duration, os.Stdin, os.Stdout)

	err = client.Connect()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("...connection error: %s\n", err))
		return
	}
	defer client.Close()

	go func() {
		err = client.Send()
		if err == nil {
			os.Stderr.WriteString("...EOF\n")
		}

		cancel()
	}()

	go func() {
		err = client.Receive()
		if err == nil {
			os.Stderr.WriteString("...Connection was closed by peer\n")
		}
	}()

	<-ctx.Done()
}
