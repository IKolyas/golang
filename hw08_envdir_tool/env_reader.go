package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// isBinaryData проверяет, содержит ли данные бинарные символы.
func isBinaryData(data []byte) bool {
	for _, b := range data {
		if b == 0x00 {
			continue
		}
		if !unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			return true
		}
	}
	return false
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		if isBinaryData(content) {
			continue
		}

		lines := bytes.Split(content, []byte("\n"))
		if len(lines) == 0 {
			continue
		}

		firstLine := lines[0]
		firstLine = bytes.ReplaceAll(firstLine, []byte{0x00}, []byte("\n"))
		re := regexp.MustCompile(`[[:cntrl:]]`)
		firstLine = re.ReplaceAll(firstLine, []byte("\n"))
		value := strings.TrimRight(string(firstLine), " \t\n")

		if len(value) == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
		} else {
			env[file.Name()] = EnvValue{Value: value}
		}
	}

	return env, nil
}
