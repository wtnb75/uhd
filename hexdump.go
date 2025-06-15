package main

import (
	"fmt"
	"io"
	"log/slog"
)

type hexdump struct {
	output io.Writer
	cur    uint64
	width  int
	sep    int
}

func (h *hexdump) Write(p []byte) (n int, err error) {
	for i, ch := range p {
		fmt.Fprintf(h.output, " %02X", uint8(ch))
		if (h.cur+uint64(i))%uint64(h.sep) == uint64(h.sep-1) {
			fmt.Fprint(h.output, " ")
		}
		if (h.cur+uint64(i))%uint64(h.width) == uint64(h.width)-1 {
			fmt.Fprint(h.output, "\n")
		}
	}
	h.cur += uint64(len(p))
	return len(p), nil
}

func (h *hexdump) Close() (err error) {
	if h.cur%uint64(h.width) != 0 {
		fmt.Fprint(h.output, "\n")
	}
	if closer, ok := h.output.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			slog.Error("close writer", "err", err)
			return err
		}
	}
	return nil
}

func NewHexdump(output io.Writer, width int, sep int) *hexdump {
	return &hexdump{
		output: output,
		cur:    0,
		width:  width,
		sep:    sep,
	}
}
