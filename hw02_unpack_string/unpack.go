package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder

	input := []rune(s)

	for i := 0; i < len(input); i++ {
		first := rune(input[i])

		if unicode.IsDigit(first) && i >= len(input)-1 {
			continue
		}

		if i == len(input)-1 && unicode.IsLetter(first) {
			result.WriteString(string(input[i]))

			continue
		}

		second := rune(input[i+1])

		if unicode.IsDigit(first) {
			if i == 0 || unicode.IsDigit(second) {
				return "", ErrInvalidString
			}

			continue
		}

		if unicode.IsLetter(second) {
			result.WriteString(string(first))

			continue
		}

		n, _ := strconv.Atoi(string(second))

		result.WriteString(strings.Repeat(string(first), n))
	}

	return result.String(), nil
}
