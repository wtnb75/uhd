package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestHexbytes(t *testing.T) {
	buf := &bytes.Buffer{}
	hb := NewHexbytes(buf, 16)
	input := "hello world 1234567890"
	expected := []int{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f,
		0x72, 0x6c, 0x64, 0x20, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30}
	n, err := fmt.Fprint(hb, input)
	if err != nil {
		t.Error("write", "err", err)
	}
	if err = hb.Close(); err != nil {
		t.Error("close", "err", err)
	}
	if n != len(input) {
		t.Error("short write", "n", n, "expected", len(input))
	}
	parsed := make([]int, 0, len(input))
	for _, token := range strings.Split(buf.String(), ",") {
		tk := strings.Trim(token, " \n")
		if tk == "" {
			continue
		}
		if val, err := strconv.ParseUint(tk, 0, 8); err != nil {
			t.Error("parseint", "err", err, "token", tk)
		} else {
			parsed = append(parsed, int(val))
		}
	}
	if !reflect.DeepEqual(parsed, expected) {
		t.Error("mismatch", "output", parsed, "expected", expected)
	}
}
