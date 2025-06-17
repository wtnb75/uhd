package main

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type hexrev struct {
	output io.Writer
}

func (h *hexrev) Write(p []byte) (n int, err error) {
	var ch byte
	for _, hexstr := range strings.Split(string(p), " ") {
		n, err := fmt.Sscanf(hexstr, "%02X", &ch)
		if err == io.EOF {
			slog.Debug("eof", "hexstr", hexstr)
			break
		}
		if n != 1 {
			slog.Info("scan failed", "n", n, "hexstr", hexstr)
			continue
		}
		if err != nil {
			slog.Error("scan", "err", err, "str", hexstr)
			return len(p), err
		}
		written, err := h.output.Write([]byte{ch})
		if err != nil {
			slog.Error("write", "err", err, "ch", ch, "written", written)
			return len(p), err
		}
		if written != 1 {
			slog.Warn("short write", "written", written, "ch", ch)
		}
	}
	return len(p), nil
}

func (h *hexrev) Close() error {
	return nil
}

func NewHexrev(output io.Writer) *hexrev {
	return &hexrev{
		output: output,
	}
}
