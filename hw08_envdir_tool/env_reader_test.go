package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDirBinaryData(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "BINARY")
	err := os.WriteFile(filePath, []byte("\x00\x01\x02"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file %s", err)
	}

	env, err := ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	if _, exists := env["BINARY"]; exists {
		t.Error("Key BINARY should have been ignored")
	}
}
