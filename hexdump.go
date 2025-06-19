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
	lower  bool
}

func (h *hexdump) Write(p []byte) (n int, err error) {
	for i, ch := range p {
		if h.lower {
			fmt.Fprintf(h.output, " %02x", uint8(ch))
		} else {
			fmt.Fprintf(h.output, " %02X", uint8(ch))
		}
		c := h.cur + uint64(i)
		cw := int(c % uint64(h.width))
		if cw == h.width-1 {
			fmt.Fprint(h.output, "\n")
		} else if cw%h.sep == h.sep-1 {
			fmt.Fprint(h.output, " ")
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
		lower:  false,
	}
}

func NewHexdumpLower(output io.Writer, width int, sep int) *hexdump {
	return &hexdump{
		output: output,
		cur:    0,
		width:  width,
		sep:    sep,
		lower:  true,
	}
}
