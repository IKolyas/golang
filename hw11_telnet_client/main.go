package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
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

	timeout, err := getTimeout()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing timeout: %v\n", err)
		os.Exit(1)
	}

	args := flag.Args()

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=<timeout>] <host> <port>")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	address := fmt.Sprintf("%s:%s", host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	// Подключение к серверу
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()
	// Канал для обработки ошибок
	errChan := make(chan error, 2)
	// Горутина для отправки данных
	go func() {
		if err := client.Send(); err != nil {
			errChan <- fmt.Errorf("error sending data: %w", err)
		}
	}()
	// Горутина для получения данных
	go func() {
		if err := client.Receive(); err != nil {
			if errors.Is(err, io.EOF) {
				// Соединение закрыто корректно
				return
			}
			errChan <- fmt.Errorf("error receiving data: %w", err)
		}
	}()
	// Канал для обработки сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	// Ожидание сигнала завершения или ошибки
	select {
	case sig := <-sigChan:
		fmt.Fprintf(os.Stderr, "...%s received, closing connection\n", sig)
	case err := <-errChan:
		fmt.Fprintln(os.Stderr, err)
	}
}

func getTimeout() (time.Duration, error) {
	var timeout string
	flag.StringVar(&timeout, "timeout", "10s", "timeout in seconds")
	flag.Parse()

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %w", err)
	}

	if duration <= 0 {
		return 10 * time.Second, nil
	}

	return duration, nil
}
