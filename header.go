package main

import (
	"fmt"
	"io"
	"log/slog"
)

type header struct {
	output io.Writer
	cur    uint64
	width  int
}

func (h *header) Write(p []byte) (n int, err error) {
	for i := h.cur; i < h.cur+uint64(len(p)); i++ {
		if i%uint64(h.width) == 0 {
			fmt.Fprintf(h.output, "%08X\n", i)
		}
	}
	h.cur += uint64(len(p))
	return len(p), nil
}

func (h *header) Close() (err error) {
	if closer, ok := h.output.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			slog.Error("close writer", "err", err)
			return err
		}
	}
	return nil
}

func NewHeader(output io.Writer, width int) *header {
	return &header{
		output: output,
		cur:    0,
		width:  width,
	}
}
