package main

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/width"
)

type printable struct {
	output   io.Writer
	cur      uint64
	width    int
	encoding string
	rest     []byte
}

func (h *printable) writeShiftJIS(p []byte) (n int, err error) {
	dec := japanese.ShiftJIS.NewDecoder()
	runesrc := make([]byte, 0, 2)
	mb := false
	for _, ch := range p {
		if len(runesrc) == 0 && ch >= 0x80 {
			runesrc = append(runesrc, ch)
			continue
		} else if len(runesrc) == 1 {
			runesrc = append(runesrc, ch)
			runesrc_u8, err := dec.Bytes(runesrc)
			if err != nil {
				fmt.Fprintf(h.output, ".")
				runesrc = runesrc[1:]
				h.cur += 1
				mb = false
			} else {
				r, _ := utf8.DecodeRune(runesrc_u8)
				if unicode.IsPrint(r) {
					charwidth := h.runeWidth(r)
					fmt.Fprintf(h.output, "%c", r)
					if charwidth == 1 {
						fmt.Fprint(h.output, "_")
					}
				} else {
					fmt.Fprint(h.output, "..")
				}
				runesrc = make([]byte, 0, 2)
				h.cur += 2
				mb = true
			}
		} else {
			if 0x20 <= ch && ch <= 0x7e {
				fmt.Fprint(h.output, string(ch))
			} else {
				fmt.Fprint(h.output, ".")
			}
			h.cur += 1
			mb = false
		}
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, "\n")
		} else if mb && h.cur%uint64(h.width) == 1 {
			fmt.Fprint(h.output, "\n_")
		}
	}
	return len(p), nil
}

func (h *printable) writeEUCJP(p []byte) (n int, err error) {
	dec := japanese.EUCJP.NewDecoder()
	runesrc := make([]byte, 0, 2)
	mb := false
	for _, ch := range p {
		if ch >= 0x80 {
			runesrc = append(runesrc, ch)
			if len(runesrc) == 2 {
				runesrc_u8, err := dec.Bytes(runesrc)
				if err != nil {
					fmt.Fprintf(h.output, "..")
				} else {
					r, _ := utf8.DecodeRune(runesrc_u8)
					if unicode.IsPrint(r) {
						charwidth := h.runeWidth(r)
						fmt.Fprintf(h.output, "%c", r)
						if charwidth == 1 {
							fmt.Fprint(h.output, "_")
						}
					} else {
						fmt.Fprint(h.output, "..")
					}
				}
				runesrc = make([]byte, 0, 2)
				h.cur += 2
				mb = true
			} else {
				continue
			}
		} else {
			if len(runesrc) > 0 {
				fmt.Fprint(h.output, ".")
				h.cur += 1
				runesrc = make([]byte, 0, 2)
				if h.cur%uint64(h.width) == 0 {
					fmt.Fprint(h.output, "\n")
				}
			}
			if 0x20 <= ch && ch <= 0x7e {
				fmt.Fprint(h.output, string(ch))
			} else {
				fmt.Fprint(h.output, ".")
			}
			h.cur += 1
			mb = false
		}
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, "\n")
		}
		if mb && h.cur%uint64(h.width) == 1 {
			fmt.Fprint(h.output, "\n_")
		}
	}
	return len(p), nil
}

func (h *printable) writeASCII(p []byte) (n int, err error) {
	for _, ch := range p {
		if 0x20 <= ch && ch <= 0x7e {
			fmt.Fprint(h.output, string(ch))
		} else {
			fmt.Fprint(h.output, ".")
		}
		h.cur += 1
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, "\n")
		}
	}
	return len(p), nil
}

func (h *printable) runeWidth(r rune) int {
	prop := width.LookupRune(r)
	if prop.Kind() == width.EastAsianNarrow || prop.Kind() == width.EastAsianHalfwidth || prop.Kind() == width.Neutral {
		return 1
	}
	return 2
}

func (h *printable) writeUTF8(p []byte) (n int, err error) {
	h.rest = append(h.rest, p...)
	for len(h.rest) > 0 {
		curpos := int(h.cur % uint64(h.width))
		r, size := utf8.DecodeRune(h.rest)
		if r == utf8.RuneError && size == 1 {
			if len(h.rest) < 4 {
				return len(p), nil
			}
			fmt.Fprint(h.output, ".")
		} else if !unicode.IsPrint(r) {
			fmt.Fprint(h.output, ".")
		} else {
			fmt.Fprintf(h.output, "%c", r)
		}
		charwidth := h.runeWidth(r)
		if curpos+size < h.width && size > charwidth {
			fmt.Fprint(h.output, strings.Repeat("_", size-charwidth))
		}
		h.rest = h.rest[size:]
		if curpos+size >= h.width {
			fmt.Fprintf(h.output, "\n")
			fill := (curpos + size) % h.width
			if fill > 0 {
				fmt.Fprint(h.output, strings.Repeat("_", fill))
			}
		}
		h.cur += uint64(size)
	}
	return len(p), nil
}

func (h *printable) Write(p []byte) (n int, err error) {
	if h.encoding == "utf-8" {
		return h.writeUTF8(p)
	} else if h.encoding == "euc-jp" {
		return h.writeEUCJP(p)
	} else if h.encoding == "shift-jis" {
		return h.writeShiftJIS(p)
	}
	return h.writeASCII(p)
}

func (h *printable) Close() (err error) {
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

func NewPrintable(output io.Writer, encoding string, width int) *printable {
	return &printable{
		output:   output,
		cur:      0,
		width:    width,
		encoding: encoding,
		rest:     make([]byte, 0),
	}
}
