package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	var count int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			count = count*10 + int(r-'0')
		} else {
			if count == 0 {
				count = 1
			}
			for i := 0; i < count; i++ {
				result.WriteRune(r)
			}
			count = 0
		}
	}
	if count != 0 {
		return "", fmt.Errorf("invalid string")
	}
	return result.String(), nil
}
