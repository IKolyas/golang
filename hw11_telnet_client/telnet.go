package main

import (
	"bufio"
	"errors"
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

type TClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *TClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", tc.address)
	return nil
}

func (tc *TClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

func (tc *TClient) Send() error {
	scanner := bufio.NewScanner(tc.in)
	for scanner.Scan() {
		_, err := tc.conn.Write(scanner.Bytes())
		if err != nil {
			return err
		}
		_, err = tc.conn.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "...EOF")
	return nil
}

func (tc *TClient) Receive() error {
	reader := bufio.NewReader(tc.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
				return nil
			}
			return err
		}
		_, err = fmt.Fprint(tc.out, line)
		if err != nil {
			return err
		}
	}
}
