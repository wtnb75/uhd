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
	start_ch string
	end_ch   string
}

type uintrange struct {
	start uint
	end   uint
}

func sjis_single(ch byte) bool {
	return ch <= 0x7e || (0xa1 <= ch && ch <= 0xdf)
}

func valid_sjis(b1, b2 byte) bool {
	// info from http://charset.7jp.net/sjis.html
	invalid := []uintrange{
		{0x81ad, 0x81b7},
		{0x81c0, 0x81c7},
		{0x81cf, 0x81d9},
		{0x81e9, 0x81ef},
		{0x81f8, 0x81fb},

		{0x8240, 0x824e},
		{0x8259, 0x825f},
		{0x827a, 0x8280},
		{0x829b, 0x829e},
		{0x82f2, 0x82fc},

		{0x8397, 0x839e},
		{0x83b7, 0x83be},
		{0x83d7, 0x83fc},

		{0x8461, 0x846f},
		{0x8492, 0x849e},
		{0x84bf, 0x84fc},

		{0x8540, 0x889e},

		{0x9873, 0x989e},

		{0xa040, 0xdffc},
		{0xeaa5, 0xeffc},
	}
	ch := (uint(b1) << 8) | uint(b2)
	for _, r := range invalid {
		if r.start <= ch && ch <= r.end {
			slog.Debug("invalid euc-jp", "b1", b1, "b2", b2, "ch", ch, "range", r)
			return false
		}
	}
	return true
}

func (h *printable) writeShiftJIS(p []byte) (n int, err error) {
	dec := japanese.ShiftJIS.NewDecoder()
	runesrc := make([]byte, 0, 2)
	mb := false
	for _, ch := range p {
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		if len(runesrc) == 0 && !sjis_single(ch) {
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
				r, size := utf8.DecodeRune(runesrc_u8)
				slog.Debug("decoded", "rune", r, "size", size)
				if !utf8.ValidRune(r) || !valid_sjis(runesrc[0], runesrc[1]) {
					fmt.Fprint(h.output, "..")
				} else if unicode.IsPrint(r) {
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
			} else if 0xa1 <= ch && ch <= 0xdf {
				runesrc_u8, err := dec.Bytes([]byte{ch})
				if err != nil {
					fmt.Fprint(h.output, ".")
				} else {
					r, _ := utf8.DecodeRune(runesrc_u8)
					fmt.Fprintf(h.output, "%c", r)
				}
			} else {
				fmt.Fprint(h.output, ".")
			}
			h.cur += 1
			mb = false
		}
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.end_ch+"\n")
		} else if mb && h.cur%uint64(h.width) == 1 {
			fmt.Fprint(h.output, "\n"+h.start_ch+"_")
		}
	}
	return len(p), nil
}

func valid_eucjp(b1, b2 byte) bool {
	// info from http://charset.7jp.net/euc.html
	invalid := []uintrange{
		{0xa2af, 0xa2b9},
		{0xa2c2, 0xa2c9},
		{0xa2d1, 0xa2db},
		{0xa2eb, 0xa2f1},
		{0xa2fa, 0xa2fd},

		{0xa3a1, 0xa3af},
		{0xa3ba, 0xa3c0},
		{0xa3db, 0xa3e0},
		{0xa3fb, 0xa3fe},

		{0xa4f4, 0xa4fe},

		{0xa5f7, 0xa5fe},

		{0xa6b9, 0xa6c0},
		{0xa6d9, 0xa6fe},

		{0xa7c2, 0xa7d0},
		{0xa7f2, 0xa7fe},

		{0xa8c1, 0xa8fe},

		{0xa9a1, 0xaffe},
		{0xf4a7, 0xfefe},
	}
	ch := (uint(b1) << 8) | uint(b2)
	for _, r := range invalid {
		if r.start <= ch && ch <= r.end {
			slog.Debug("invalid euc-jp", "b1", b1, "b2", b2, "ch", ch, "range", r)
			return false
		}
	}
	return true
}

func (h *printable) writeEUCJP(p []byte) (n int, err error) {
	dec := japanese.EUCJP.NewDecoder()
	runesrc := make([]byte, 0, 2)
	mb := false
	for _, ch := range p {
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		if (0xa1 <= ch && ch <= 0xfe) || ch == 0x8e {
			runesrc = append(runesrc, ch)
			if len(runesrc) == 2 {
				runesrc_u8, err := dec.Bytes(runesrc)
				if err != nil {
					fmt.Fprintf(h.output, "..")
				} else {
					r, _ := utf8.DecodeRune(runesrc_u8)
					if !utf8.ValidRune(r) || !valid_eucjp(runesrc[0], runesrc[1]) {
						fmt.Fprint(h.output, "..")
					} else if unicode.IsPrint(r) {
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
					fmt.Fprint(h.output, h.end_ch+"\n")
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
			fmt.Fprint(h.output, h.end_ch+"\n")
		}
		if mb && h.cur%uint64(h.width) == 1 {
			fmt.Fprint(h.output, "\n"+h.start_ch+"_")
		}
	}
	return len(p), nil
}

func (h *printable) writeASCII(p []byte) (n int, err error) {
	for _, ch := range p {
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		if 0x20 <= ch && ch <= 0x7e {
			fmt.Fprint(h.output, string(ch))
		} else {
			fmt.Fprint(h.output, ".")
		}
		h.cur += 1
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.end_ch+"\n")
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
		if curpos == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
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
		if curpos+size <= h.width && size > charwidth {
			fmt.Fprint(h.output, strings.Repeat("_", size-charwidth))
		}
		h.rest = h.rest[size:]
		if curpos+size >= h.width {
			if curpos+size == h.width {
				fmt.Fprint(h.output, h.end_ch)
			}
			fmt.Fprint(h.output, "\n")
			fill := (curpos + size) % h.width
			if fill > 0 {
				fmt.Fprint(h.output, h.start_ch+strings.Repeat("_", fill))
			}
		}
		h.cur += uint64(size)
	}
	return len(p), nil
}

func (h *printable) Write(p []byte) (n int, err error) {
	switch strings.ToLower(h.encoding) {
	case "utf-8", "utf8":
		return h.writeUTF8(p)
	case "euc-jp", "eucjp":
		return h.writeEUCJP(p)
	case "shift-jis", "sjis", "shiftjis":
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

func NewPrintableSep(output io.Writer, encoding string, width int, start_ch, end_ch string) *printable {
	return &printable{
		output:   output,
		cur:      0,
		width:    width,
		encoding: encoding,
		rest:     make([]byte, 0),
		start_ch: start_ch,
		end_ch:   end_ch,
	}
}
