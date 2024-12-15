package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder

	for i := 0; i < len(s); i++ {

		if !unicode.IsDigit(rune(s[i])) {
			result.WriteString(string(s[i]))
			continue
		}

		if i == 0 {
			return "", fmt.Errorf("invalid string")
		}

		if i+1 <= len(s) && unicode.IsDigit(rune(s[i+1])) {
			return "", fmt.Errorf("invalid string")
		}

		count, _ := strconv.Atoi(string(s[i]))
		if count == 0 {
			continue
		}

		for n := count; n > 1; n-- {
			result.WriteString(string(s[i-1]))
		}
	}
	return result.String(), nil
}
