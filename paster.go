package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
)

type paster struct {
	writer  io.Writer
	dupidx  int
	readers []*bufio.Scanner
}

func (p *paster) Process(widths ...int) error {
	var prev []string
	var in_dup = false
	txts := make([]string, 0, len(p.readers))
	for {
		eof := true
		for _, reader := range p.readers {
			got := reader.Scan()
			if !got {
				continue
			}
			eof = false
			txt := reader.Text()
			txts = append(txts, txt)
			if reader.Err() != nil {
				slog.Error("scan error", "err", reader.Err())
				return reader.Err()
			}
		}
		if len(txts) == 0 {
			break
		}
		if len(prev) != 0 && txts[p.dupidx] == prev[p.dupidx] {
			if !in_dup {
				fmt.Fprint(p.writer, "*\n")
			}
			in_dup = true
		} else {
			in_dup = false
			for idx, txt := range txts {
				fmt.Fprintf(p.writer, "%-*s", widths[idx], txt)
			}
			if eof {
				slog.Debug("eof")
				break
			}
			fmt.Fprint(p.writer, "\n")
		}
		prev = txts[:]
		txts = make([]string, 0, len(p.readers))
		slog.Debug("line", "prev", prev, "txts", txts)
	}
	return nil
}

func NewPaster(writer io.Writer, dupidx int, readers ...io.Reader) *paster {
	rds := make([]*bufio.Scanner, 0)
	for _, rd := range readers {
		rds = append(rds, bufio.NewScanner(rd))
	}
	return &paster{
		writer:  writer,
		dupidx:  dupidx,
		readers: rds,
	}
}
