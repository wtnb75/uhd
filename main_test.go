package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldOption := option
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		option = oldOption
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("pipe", err)
	}
	os.Stdout = w
	os.Args = []string{"uhd", "--version"}

	main()

	if err := w.Close(); err != nil {
		t.Error("write close", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Error("read", err)
	}

	if !strings.Contains(string(out), "uhd") {
		t.Error("missing command name", string(out))
	}
	if !strings.Contains(string(out), version) {
		t.Error("missing version", string(out))
	}
	if !strings.Contains(string(out), commit) {
		t.Error("missing commit", string(out))
	}
	if !strings.Contains(string(out), date) {
		t.Error("missing date", string(out))
	}
}
