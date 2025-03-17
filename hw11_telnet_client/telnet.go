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

	in io.ReadCloser

	out io.Writer

	conn net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {

	return &TClient{

		address: address,

		timeout: timeout,

		in: in,

		out: out,
	}

}

func (tc *TClient) Connect() error {

	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)

	if err != nil {

		return fmt.Errorf("ошибка подключения к %s: %w", tc.address, err)

	}

	tc.conn = conn

	fmt.Fprintf(os.Stderr, "...Подключено к %s\n", tc.address)

	return nil

}

func (tc *TClient) Close() error {

	if tc.conn != nil {

		if err := tc.conn.Close(); err != nil {

			return fmt.Errorf("ошибка закрытия соединения: %w", err)

		}

	}

	return nil

}

func (tc *TClient) Send() error {

	scanner := bufio.NewScanner(tc.in)

	for scanner.Scan() {

		input := scanner.Text()

		_, err := fmt.Fprintf(tc.conn, "%s\n", input)

		if err != nil {

			return fmt.Errorf("ошибка отправки данных: %w", err)

		}

	}

	if err := scanner.Err(); err != nil {

		return fmt.Errorf("ошибка чтения ввода: %w", err)

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

				fmt.Fprintln(os.Stderr, "...Соединение закрыто удаленной стороной")

				return nil

			}

			return fmt.Errorf("ошибка получения данных: %w", err)

		}

		_, err = fmt.Fprint(tc.out, line)

		if err != nil {

			return fmt.Errorf("ошибка вывода данных: %w", err)

		}

	}

}
