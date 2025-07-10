package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/acomagu/bufpipe"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"golang.org/x/text/encoding/charmap"
)

var option struct {
	Verbose  bool   `short:"v" long:"verbose" description:"Enable verbose logging"`
	Encoding string `long:"encoding" default:"utf-8"`
	Width    int    `long:"width" default:"16"`
	Sep      int    `long:"sep" default:"8"`
	Layout   string `long:"layout" default:"jhd" choice:"hexdump" choice:"jhd" choice:"bytes"`
	ListCode bool   `short:"l" long:"list-codes" description:"list encoding"`
	NoColor  bool   `long:"no-color" description:"disable color output"`
}

type column struct {
	name  string
	width int
}

func get_layout(predefined string) []column {
	switch predefined {
	case "jhd":
		return []column{
			{"header", 9},
			{"hexdump", 3*(option.Width) + option.Width/option.Sep + (option.Width / 8) + 1},
			{"printable", option.Width},
		}
	case "hexdump":
		return []column{
			{"header", 9},
			{"hexdump_lower", 3*(option.Width) + option.Width/option.Sep + (option.Width / 8) + 1},
			{"printable_pipe", option.Width + 2},
		}
	case "bytes":
		return []column{
			{"header", 9},
			{"hexbytes_lower", 6*option.Width + 1},
			{"printable", option.Width},
		}
	}
	return []column{}
}

func do_uhd(filename string) (err error) {
	var rd *os.File
	var layout = get_layout(option.Layout)
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
	widths := make([]int, 0, len(layout))
	writers := make([]io.Writer, 0, len(layout))
	readers := make([]io.Reader, 0, len(layout))
	var dupidx int
	for idx, col := range layout {
		r, w := bufpipe.New(nil)
		switch col.name {
		case "header":
			writers = append(writers, NewHeader(w, option.Width))
		case "header_lower":
			writers = append(writers, NewHeaderLower(w, option.Width))
		case "hexdump":
			writers = append(writers, NewHexdump(w, option.Width, option.Sep))
			dupidx = idx
		case "hexdump_lower":
			writers = append(writers, NewHexdumpLower(w, option.Width, option.Sep))
			dupidx = idx
		case "hexbytes":
			writers = append(writers, NewHexbytes(w, option.Width))
		case "hexbytes_lower":
			writers = append(writers, NewHexbytesLower(w, option.Width))
		case "printable":
			writers = append(writers, NewPrintable(w, option.Encoding, option.Width))
		case "printable_pipe":
			writers = append(writers, NewPrintableSep(w, option.Encoding, option.Width, "|", "|"))
		}
		readers = append(readers, r)
		widths = append(widths, col.width)
	}
	wr := io.MultiWriter(writers...)
	pst := NewPaster(os.Stdout, dupidx, readers...)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		slog.Debug("widths", "values", widths)
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
		return
	}
	if option.NoColor {
		color.NoColor = true
	}
	if option.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if option.ListCode {
		fmt.Println("utf-8, utf8")
		fmt.Println("utf-16, utf16, utf-16be, utf16be, utf-16le, utf16le")
		fmt.Println("utf-32, utf32, utf-32be, utf32be, utf-23le, utf32le")
		fmt.Println("euc-jp, eucjp")
		fmt.Println("euc-kr, euckr")
		fmt.Println("euc-cn, euccn, gb18030")
		fmt.Println("big5")
		fmt.Println("shift-jis, sjis, shiftjis, cp932, cp-932, windows-31j")
		for _, cm := range charmap.All {
			name := fmt.Sprintf("%s", cm)
			if strings.Contains(name, "enc=") {
				tok := strings.SplitN(name, "enc=", 2)
				if len(tok) == 2 {
					name = strings.Trim(tok[1], "\"")
				}
			}
			fmt.Println(name)
		}
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
