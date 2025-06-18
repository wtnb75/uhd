package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHexdump(t *testing.T) {
	buf := &bytes.Buffer{}
	hex := NewHexdump(buf, 16, 8)
	input := "hello world 1234567890"
	expected := " 68 65 6C 6C 6F 20 77 6F  72 6C 64 20 31 32 33 34\n 35 36 37 38 39 30\n"
	n, err := fmt.Fprint(hex, input)
	if err != nil {
		t.Error("write", "err", err)
	}
	if err = hex.Close(); err != nil {
		t.Error("close", "err", err)
	}
	if n != len(input) {
		t.Error("short write", "n", n, "expected", len(input))
	}
	output := buf.String()
	if output != expected {
		t.Error("mismatch", "output", output, "expected", expected)
	}
}
