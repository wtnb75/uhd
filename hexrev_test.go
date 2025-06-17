package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHexrev(t *testing.T) {
	textdata := "01 12 34 56 78 9a bc de f0\n"
	buf := &bytes.Buffer{}
	hr := NewHexrev(buf)
	n, err := fmt.Fprint(hr, textdata)
	if err != nil {
		t.Error("fprint", "err", err)
	}
	if err = hr.Close(); err != nil {
		t.Error("close", "err", err)
	}
	if n != len(textdata) {
		t.Error("short write", "n", n, "len(textdata)", len(textdata))
	}
	if !bytes.Equal(buf.Bytes(), []byte{0x01, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}) {
		t.Error("mismatch", "buf.Bytes()", buf.Bytes(), "expected", textdata)
	}
}
