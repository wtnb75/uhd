package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
)

type paster struct {
	writer  io.Writer
	readers []*bufio.Scanner
}

func (p *paster) Process(widths ...int) error {
	for {
		eof := true
		for idx, reader := range p.readers {
			got := reader.Scan()
			if !got {
				continue
			}
			eof = false
			fmt.Fprintf(p.writer, "%-*s", widths[idx], reader.Text())
			if reader.Err() != nil {
				slog.Error("scan error", "err", reader.Err())
				return reader.Err()
			}
		}
		if eof {
			slog.Debug("eof")
			break
		}
		fmt.Fprint(p.writer, "\n")
	}
	return nil
}

func NewPaster(writer io.Writer, readers ...io.Reader) *paster {
	rds := make([]*bufio.Scanner, 0)
	for _, rd := range readers {
		rds = append(rds, bufio.NewScanner(rd))
	}
	return &paster{
		writer:  writer,
		readers: rds,
	}
}
