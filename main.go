package main

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/acomagu/bufpipe"
	"github.com/jessevdk/go-flags"
)

var option struct {
	Verbose  bool   `short:"v" long:"verbose" description:"Enable verbose logging"`
	Encoding string `long:"encoding" default:"utf-8"`
	Width    int    `long:"width" default:"16"`
	Sep      int    `long:"sep" default:"8"`
}

type column struct {
	name  string
	width int
}

func do_uhd(filename string) (err error) {
	layout := []column{
		column{"header", 9},
		column{"hexdump", 3*(option.Width) + (option.Width / 8) + 2},
		column{"printable", option.Width},
	}
	var rd *os.File
	if filename == "-" {
		rd = os.Stdin
	} else {
		rd, err = os.Open(filename)
		if err != nil {
			slog.Error("open", "file", filename, "err", err)
			return err
		}
		defer rd.Close()
	}
	widths := make([]int, 0)
	writers := make([]io.Writer, 0)
	readers := make([]io.Reader, 0)
	for _, col := range layout {
		r, w := bufpipe.New(nil)
		if col.name == "header" {
			writers = append(writers, NewHeader(w, option.Width))
		} else if col.name == "hexdump" {
			writers = append(writers, NewHexdump(w, option.Width, option.Sep))
		} else if col.name == "printable" {
			writers = append(writers, NewPrintable(w, option.Encoding, option.Width))
		}
		readers = append(readers, r)
		widths = append(widths, col.width)
	}
	wr := io.MultiWriter(writers...)
	pst := NewPaster(os.Stdout, readers...)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := pst.Process(widths...); err != nil {
			slog.Error("paster", "err", err)
		}
		slog.Debug("finished", "file", filename)
		wg.Done()
	}()
	written, err := io.Copy(wr, rd)
	slog.Debug("copy", "file", filename, "written", written, "err", err)
	if err != nil {
		slog.Error("copy", "file", filename, "err", err)
	}
	for _, w := range writers {
		if closer, ok := w.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				slog.Error("close writer", "file", filename, "err", err)
			}
		}
	}
	wg.Wait()
	return nil
}

func main() {
	parser := flags.NewParser(&option, flags.Default)
	parsed, err := parser.Parse()
	if err != nil {
		slog.Error("parse", "err", err)
		return
	}
	if option.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if option.Width%option.Sep != 0 {
		slog.Error("width must be a multiple of sep", "width", option.Width, "sep", option.Sep)
		return
	}
	if len(parsed) == 0 {
		err := do_uhd("-")
		if err != nil {
			slog.Error("uhd", "file", "(stdin)", "err", err)
		}
	} else {
		for _, fn := range parsed {
			err := do_uhd(fn)
			if err != nil {
				slog.Error("uhd", "file", fn, "err", err)
				// continue
			}
		}
	}
}
