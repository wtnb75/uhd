package main

import (
	"fmt"
	"io"
	"log/slog"
)

type hexbytes struct {
	output io.Writer
	cur    uint64
	width  int
	lower  bool
}

func (h *hexbytes) Write(p []byte) (n int, err error) {
	for i, ch := range p {
		if h.lower {
			fmt.Fprintf(h.output, "0x%02x,", uint8(ch))
		} else {
			fmt.Fprintf(h.output, "0x%02X,", uint8(ch))
		}
		if i%h.width == h.width-1 {
			fmt.Fprint(h.output, "\n")
		} else {
			fmt.Fprint(h.output, " ")
		}
	}
	h.cur += uint64(len(p))
	return len(p), nil
}

func (h *hexbytes) Close() (err error) {
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

func NewHexbytes(output io.Writer, width int) *hexbytes {
	return &hexbytes{
		output: output,
		cur:    0,
		width:  width,
		lower:  false,
	}
}

func NewHexbytesLower(output io.Writer, width int) *hexbytes {
	return &hexbytes{
		output: output,
		cur:    0,
		width:  width,
		lower:  true,
	}
}
