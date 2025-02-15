package main

import (
	"os"
)

func main() {
	if len(os.Args) < 3 {
		os.Stderr.WriteString("Usage: go-envdir /path/to/env/dir command [args...]\n")
		os.Exit(1)
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		os.Stderr.WriteString("Error reading environment directory: " + err.Error() + "\n")
		os.Exit(1)
	}

	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
