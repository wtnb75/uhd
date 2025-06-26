package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
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
	lendian  bool
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
	p = append(h.rest, p...)
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
	h.rest = runesrc
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

func valid_euckr(b1, b2 byte) bool {
	// info from https://uic.jp/charset/show/euc-kr/
	invalid := []uintrange{
		{0xa2e8, 0xa2ff},
		{0xa5ab, 0xa5af},
		{0xa5ba, 0xa5c0},
		{0xa5d9, 0xa5e0},
		{0xa5f9, 0xa6a0},
		{0xa6e5, 0xa7a0},
		{0xa8a5, 0xa8a5},
		{0xa8a7, 0xa8a7},
		{0xa8b0, 0xa8b0},
		{0xaaf4, 0xaaff},
		{0xabf7, 0xabff},
		{0xacc2, 0xacd0},
		{0xacf2, 0xacff},
		{0xfdff, 0xffff},
	}
	if b2 == 0xa0 || b2 == 0xff {
		return false
	}
	ch := (uint(b1) << 8) | uint(b2)
	for _, r := range invalid {
		if r.start <= ch && ch <= r.end {
			slog.Debug("invalid euc-kr", "b1", b1, "b2", b2, "ch", ch, "range", r)
			return false
		}
	}
	return true
}

func valid_euccn(b1, b2 byte) bool {
	// info from https://uic.jp/charset/show/euc-cn/
	invalid := []uintrange{
		{0xa2e3, 0xa2e4},
		{0xa2ef, 0xa2f0},
		{0xa2fd, 0xa2ff},
		{0xa4f4, 0xa4ff},
		{0xa5f7, 0xa5ff},
		{0xa6b9, 0xa6c0},
		{0xa6d9, 0xa7a0},
		{0xa7c2, 0xa7d0},
		{0xa7f2, 0xa8a0},
		{0xa8bb, 0xa8c4},
		{0xa8ea, 0xa9a3},
		{0xa9f0, 0xb0a0},
		{0xf7ff, 0xffff},
	}
	if b2 == 0xa0 || b2 == 0xff {
		return false
	}
	ch := (uint(b1) << 8) | uint(b2)
	for _, r := range invalid {
		if r.start <= ch && ch <= r.end {
			slog.Debug("invalid euc-cn", "b1", b1, "b2", b2, "ch", ch, "range", r)
			return false
		}
	}
	return true
}

func (h *printable) writeEUCAny(p []byte, dec *encoding.Decoder, valid func(b1, b2 byte) bool) (n int, err error) {
	p = append(h.rest, p...)
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
					if !utf8.ValidRune(r) || !valid(runesrc[0], runesrc[1]) {
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
	h.rest = runesrc
	return len(p), nil
}

func (h *printable) writeEUCJP(p []byte) (n int, err error) {
	dec := japanese.EUCJP.NewDecoder()
	return h.writeEUCAny(p, dec, valid_eucjp)
}

func (h *printable) writeEUCKR(p []byte) (n int, err error) {
	dec := korean.EUCKR.NewDecoder()
	return h.writeEUCAny(p, dec, valid_euckr)
}

func (h *printable) writeEUCCN(p []byte) (n int, err error) {
	dec := simplifiedchinese.GB18030.NewDecoder()
	return h.writeEUCAny(p, dec, valid_euccn)
}

func valid_big5(b1, b2 byte) bool {
	return true
}

func (h *printable) writeBig5(p []byte) (n int, err error) {
	dec := traditionalchinese.Big5.NewDecoder()
	runesrc := make([]byte, 0, 2)
	p = append(h.rest, p...)
	mb := false
	for _, ch := range p {
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		if len(runesrc) == 0 && 0xa1 <= ch && ch <= 0xf9 {
			runesrc = append(runesrc, ch)
			continue
		} else if len(runesrc) == 1 && 0x40 <= ch && ch <= 0xfe {
			runesrc = append(runesrc, ch)
			runesrc_u8, err := dec.Bytes(runesrc)
			if err != nil {
				fmt.Fprintf(h.output, "..")
			} else {
				r, _ := utf8.DecodeRune(runesrc_u8)
				if !utf8.ValidRune(r) || !valid_big5(runesrc[0], runesrc[1]) {
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
	h.rest = runesrc
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
	switch prop.Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	default:
		return 1
	}
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

func getcode_utf16(p []byte, lendian bool) (uint32, error) {
	if len(p) < 2 {
		return 0, fmt.Errorf("short bytes: %d", len(p))
	}
	if lendian {
		return uint32(p[1])<<8 | uint32(p[0]), nil
	}
	return uint32(p[0])<<8 | uint32(p[1]), nil
}

func (h *printable) writeUTF16(p []byte) (n int, err error) {
	p = append(h.rest, p...)
	// check bom
	cur := 0
	if p[0] == 0xff && p[1] == 0xfe {
		h.lendian = true
		fmt.Fprint(h.output, h.start_ch+"LE")
		cur = 2
	} else if p[0] == 0xfe && p[1] == 0xff {
		h.lendian = false
		fmt.Fprint(h.output, h.start_ch+"BE")
		cur = 2
	}
	for cur < len(p) {
		skip := 0
		pos := int((h.cur + uint64(cur)) % uint64(h.width))
		if pos == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		code, err := getcode_utf16(p[cur:], h.lendian)
		if err != nil {
			break
		}
		if code&0b1111_1100_0000_0000 == 0b1101_1000_0000_0000 {
			// surrogate pair 1st?
			code2, err := getcode_utf16(p[cur+2:], h.lendian)
			if err != nil {
				break
			} else if code2&0b1111_1100_0000_0000 == 0b1101_1100_0000_0000 {
				// surrogate pair 2nd
				rcode := 0x10000 + ((code & 0b0000_0011_1111_1111) << 10) | (code2 & 0b0000_0011_1111_1111)
				slog.Debug("rune", "code1", code, "code2", code2, "rune", rcode)
				fmt.Fprint(h.output, string(rune(rcode)))
				if pos+4 < h.width {
					fmt.Fprint(h.output, "__")
				}
				skip = 4
			} else {
				fmt.Fprint(h.output, ".")
				skip = 1
			}
		} else if code < 0x10000 {
			ch := rune(code)
			if unicode.IsPrint(ch) {
				charwidth := h.runeWidth(ch)
				if charwidth == 1 && pos+1 < h.width {
					fmt.Fprint(h.output, string(ch)+"_")
				} else {
					fmt.Fprint(h.output, string(ch))
				}
			} else {
				fmt.Fprint(h.output, "..")
			}
			skip = 2
		}
		if pos+skip >= h.width {
			fmt.Fprint(h.output, h.end_ch+"\n")
		}
		if pos+skip > h.width {
			fmt.Fprint(h.output, h.start_ch+strings.Repeat("_", pos+skip-h.width))
		}
		cur += skip
	}
	h.cur += uint64(cur)
	h.rest = p[cur:]
	return len(p), nil
}

func getcode_utf32(p []byte, lendian bool) (uint32, error) {
	if len(p) < 4 {
		return 0, fmt.Errorf("short bytes: %d", len(p))
	}
	if lendian {
		return uint32(p[3])<<24 | uint32(p[2])<<16 | uint32(p[1])<<8 | uint32(p[0]), nil
	}
	return uint32(p[0])<<24 | uint32(p[1])<<16 | uint32(p[2])<<8 | uint32(p[3]), nil
}

func (h *printable) writeUTF32(p []byte) (n int, err error) {
	p = append(h.rest, p...)
	// check bom
	cur := 0
	if bytes.Equal(p[:4], []byte{0x00, 0x00, 0xfe, 0xff}) {
		h.lendian = false
		fmt.Fprint(h.output, h.start_ch+"_BE_")
		cur = 4
	} else if bytes.Equal(p[:4], []byte{0xff, 0xfe, 0x00, 0x00}) {
		h.lendian = true
		fmt.Fprint(h.output, h.start_ch+"_LE_")
		cur = 4
	}
	for cur < len(p) {
		skip := 0
		pos := int((h.cur + uint64(cur)) % uint64(h.width))
		if pos == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		code, err := getcode_utf32(p[cur:], h.lendian)
		if err != nil {
			break
		}
		if code <= 0x10ffff {
			ch := rune(code)
			if unicode.IsPrint(ch) {
				charwidth := h.runeWidth(ch)
				if pos+charwidth < h.width {
					fmt.Fprint(h.output, string(ch)+strings.Repeat("_", 4-charwidth))
				} else {
					fmt.Fprint(h.output, string(ch)+strings.Repeat("_", h.width-(pos+charwidth)))
				}
			} else {
				fmt.Fprint(h.output, ".___")
			}
			skip = 4
		} else {
			skip = 1
		}
		if pos+skip >= h.width {
			fmt.Fprint(h.output, h.end_ch+"\n")
		}
		if pos+skip > h.width {
			fmt.Fprint(h.output, h.start_ch+strings.Repeat("_", pos+skip-h.width))
		}
		cur += skip
	}
	h.cur += uint64(cur)
	h.rest = p[cur:]
	return len(p), nil
}

func (h *printable) writeAny(p []byte, dec *encoding.Decoder, valid func(b []byte) bool) (n int, err error) {
	runesrc := make([]byte, 0, 2)
	runesrc = append(h.rest, runesrc...)
	mb := false
	for _, ch := range p {
		skip := 0
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.start_ch)
		}
		runesrc = append(runesrc, ch)
		runesrc_u8, err := dec.Bytes(runesrc)
		if err != nil && len(runesrc) <= 3 {
			continue
		}
		if err != nil {
			fmt.Fprintf(h.output, "..")
			runesrc = runesrc[1:]
			skip = 1
		} else {
			r, _ := utf8.DecodeRune(runesrc_u8)
			skip = len(runesrc_u8)
			runesrc = make([]byte, 0, 2)
			if !utf8.ValidRune(r) || !valid(runesrc_u8) {
				fmt.Fprint(h.output, strings.Repeat(".", len(runesrc_u8)))
			} else if unicode.IsPrint(r) {
				charwidth := h.runeWidth(r)
				fmt.Fprintf(h.output, "%c", r)
				if charwidth == 1 {
					fmt.Fprint(h.output, strings.Repeat("_", len(runesrc_u8)-charwidth))
				}
			} else {
				fmt.Fprint(h.output, strings.Repeat(".", len(runesrc_u8)))
			}
		}
		h.cur += uint64(skip)
		if h.cur%uint64(h.width) == 0 {
			fmt.Fprint(h.output, h.end_ch+"\n")
		}
		if mb && h.cur%uint64(h.width) == 1 {
			fmt.Fprint(h.output, "\n"+h.start_ch+"_")
		}
	}
	h.rest = runesrc
	return len(p), nil
}

func (h *printable) Write(p []byte) (n int, err error) {
	switch strings.ToLower(h.encoding) {
	case "utf-8", "utf8":
		return h.writeUTF8(p)
	case "utf-16", "utf16":
		return h.writeUTF16(p)
	case "utf-16be", "utf16be":
		h.lendian = false
		return h.writeUTF16(p)
	case "utf-16le", "utf16le":
		h.lendian = true
		return h.writeUTF16(p)
	case "utf-32", "utf32":
		return h.writeUTF32(p)
	case "utf-32be", "utf32be":
		h.lendian = false
		return h.writeUTF32(p)
	case "utf-23le", "utf32le":
		h.lendian = true
		return h.writeUTF32(p)
	case "euc-jp", "eucjp":
		return h.writeEUCJP(p)
	case "euc-kr", "euckr":
		return h.writeEUCKR(p)
	case "euc-cn", "euccn", "gb18030":
		return h.writeEUCCN(p)
	case "big5":
		return h.writeBig5(p)
	case "shift-jis", "sjis", "shiftjis", "cp932", "cp-932", "windows-31j":
		return h.writeShiftJIS(p)
	}
	for _, cm := range charmap.All {
		name := fmt.Sprintf("%s", cm)
		if strings.Contains(name, "enc=") {
			tok := strings.SplitN(name, "enc=", 2)
			if len(tok) == 2 {
				name = strings.Trim(tok[1], "\"")
			}
		}
		if strings.EqualFold(name, h.encoding) {
			dec := cm.NewDecoder()
			slog.Debug("using decoder", "name", name)
			return h.writeAny(p, dec, func(b []byte) bool { return true })
		}
	}
	slog.Debug("using ascii")
	return h.writeASCII(p)
}

func (h *printable) Close() (err error) {
	fmt.Fprint(h.output, strings.Repeat(".", len(h.rest)))
	h.cur += uint64(len(h.rest))
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
		lendian:  false,
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
		lendian:  false,
	}
}
