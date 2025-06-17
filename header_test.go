package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	hdr := NewHeader(buf, 16)
	n, err := hdr.Write(make([]byte, 128))
	if err != nil {
		t.Error("write", "err", err)
	}
	if err = hdr.Close(); err != nil {
		t.Error("close", "err", err)
	}
	if n != 128 {
		t.Error("short write", "n", n, "expected", 128)
	}
	rd := bufio.NewReader(buf)
	for i := 0; i < 128/16; i++ {
		line, _, err := rd.ReadLine()
		if err != nil {
			t.Error("readline", "err", err)
		}
		if string(line) != fmt.Sprintf("%08X", i*16) {
			t.Error("invalid line", "line", string(line), "i", i)
		}
	}
	if _, _, err = rd.ReadLine(); err != io.EOF {
		t.Error("no eof")
	}
}
