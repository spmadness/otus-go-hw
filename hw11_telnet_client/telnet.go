package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func (c *Client) Connect() error {
	var err error

	c.Conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	_, err = os.Stderr.WriteString(fmt.Sprintf("...Connected to %s\n", c.address))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Send() error {
	return c.scanLoop(c.In, c.Conn)
}

func (c *Client) Receive() error {
	return c.scanLoop(c.Conn, c.Out)
}

func (c *Client) scanLoop(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		str := fmt.Sprintf("%s\n", s.Text())
		_, err := w.Write([]byte(str))
		if err != nil {
			return err
		}
	}

	return s.Err()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address, timeout, in, out, nil,
	}
}
