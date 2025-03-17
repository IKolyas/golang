package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=<timeout>] <host> <port>")
		os.Exit(1)
	}
	timeout := 10 * time.Second
	args := os.Args[1:]
	// Проверка аргументов на наличие таймаута
	if len(args) > 0 && args[0][:10] == "--timeout=" {
		var err error
		timeout, err = time.ParseDuration(args[0][10:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid timeout value: %v\n", err)
			os.Exit(1)
		}
		args = args[1:]
	}
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=<timeout>] <host> <port>")
		os.Exit(1)
	}
	address := fmt.Sprintf("%s:%s", args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()
	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending data: %v\n", err)
			os.Exit(1)
		}
	}()
	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "Error receiving data: %v\n", err)
			os.Exit(1)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan
	fmt.Fprintln(os.Stderr, "...SIGINT received, closing connection")
}
