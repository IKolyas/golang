package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input string

		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},

		{input: "abccd", expected: "abccd"},

		{input: "", expected: ""},

		{input: "aaa0bsaг7", expected: "aabsaггггггг"},

		// uncomment if task with asterisk completed

		// {input: `qwe\4\5`, expected: `qwe45`},

		// {input: `qwe\45`, expected: `qwe44444`},

		// {input: `qwe\\5`, expected: `qwe\\\\\`},

		// {input: `qwe\\\3`, expected: `qwe\3`},

	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)

			require.NoError(t, err)

			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackCompareSum(t *testing.T) {
	tests := []string{"a4bc2d5e2q", "a4f5", "aaa1ab"}

	for _, tc := range tests {
		tc := tc

		t.Run(tc, func(t *testing.T) {
			sum := len(tc)

			re := regexp.MustCompile(`\d`).FindAllString(tc, -1)

			for _, num := range re {
				n, _ := strconv.Atoi(num)

				sum += n - 2
			}

			res, _ := Unpack(tc)

			require.Equal(t, len(res), sum)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}

	for _, tc := range invalidStrings {
		tc := tc

		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)

			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
