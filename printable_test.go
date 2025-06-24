package main

import (
	"bytes"
	"testing"
)

func TestPrintable_WriteASCII(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "ascii", 8)
	input := []byte("Hello,\x01!")
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "Hello,.!\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF8(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-8", 8)
	input := []byte("こんにちは世界\x00abc\x01\x7f!")
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "こ_ん_に\n_ち_は_世\n__界_.ab\nc..!\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF8_kr(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-8", 8)
	input := []byte("안녕히계십시오")
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "안_녕_히\n_계_십_시\n__오_\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF8_emoji(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-8", 16)
	input := []byte("💩や🍺などの絵文字")
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "💩__や_🍺__な_ど\n_の_絵_文_字_\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF8_hankana(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-8", 16)
	input := []byte("ﾊﾝｶｸｶﾅﾓｼﾞ")
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "ﾊ__ﾝ__ｶ__ｸ__ｶ__ﾅ\n__ﾓ__ｼ__ﾞ__\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteEUCJP(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "euc-jp", 8)
	input := []byte{0xA4, 0xB3, 0xA4, 0xF3, 0xA4, 0xCB, 0xA4, 0xC1, 0xA4, 0xCF, 0xC0, 0xA4, 0xB3, 0xA6, 0x0A, 0x61, 0x62, 0x63, 0x01, 0x7F, 0x21, 0x0A}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "こんにち\nは世界.a\nbc..!.\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteEUCJP_hankaku(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "euc-jp", 8)
	input := []byte{0x8e, 0xca, 0x8e, 0xdd, 0x8e, 0xb6, 0x8e, 0xb8, 0x8e, 0xb6, 0x8e, 0xc5, 0x8e, 0xd3, 0x8e, 0xbc, 0x8e, 0xde}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "ﾊ_ﾝ_ｶ_ｸ_\nｶ_ﾅ_ﾓ_ｼ_\nﾞ_\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteShiftJIS(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "shift-jis", 8)
	input := []byte{0x82, 0xb1, 0x82, 0xf1, 0x82, 0xc9, 0x82, 0xbf, 0x82, 0xcd, 0x90, 0xa2, 0x8a, 0x45, 0x00, 0x61, 0x62, 0x63, 0x01, 0x7f, 0x21, 0x0a}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "こんにち\nは世界.a\nbc....\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteShiftJIS_hankaku(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "shift-jis", 8)
	input := []byte{0xca, 0xdd, 0xb6, 0xb8, 0xb6, 0xc5, 0xd3, 0xbc, 0xde}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "ﾊﾝｶｸｶﾅﾓｼ\nﾞ\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF16(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf16", 8)
	input := []byte{0xfe, 0xff, 0x30, 0x53, 0x30, 0x93, 0x30, 0x6b, 0x30, 0x61, 0x30, 0x6f, 0x4e, 0x16, 0x75, 0x4c,
		0x00, 0x00, 0x00, 0x61, 0x00, 0x62, 0x00, 0x63, 0x00, 0x01, 0x00, 0x7f, 0x00, 0x21, 0x00, 0x0a}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "BEこんに\nちは世界\n..a_b_c_\n....!_..\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF16_emoji(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-16", 16)
	input := []byte{0xfe, 0xff, 0xd8, 0x3d, 0xdc, 0xa9, 0x30, 0x84, 0xd8, 0x3c, 0xdf, 0x7a, 0x30, 0x6a, 0x30, 0x69,
		0x30, 0x6e, 0x7d, 0x75, 0x65, 0x87, 0x5b, 0x57}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "BE💩__や🍺__など\nの絵文字\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF16_hankaku(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf16", 8)
	input := []byte{0xfe, 0xff, 0xff, 0x8a, 0xff, 0x9d, 0xff, 0x76, 0xff, 0x78, 0xff, 0x76, 0xff, 0x85, 0xff, 0x93,
		0xff, 0x7c, 0xff, 0x9e}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "BEﾊ_ﾝ_ｶ_\nｸ_ｶ_ﾅ_ﾓ_\nｼ_ﾞ_\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF32(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf32", 8)
	input := []byte{0x00, 0x00, 0xfe, 0xff, 0x00, 0x00, 0x30, 0x53, 0x00, 0x00, 0x30, 0x93, 0x00, 0x00, 0x30, 0x6b,
		0x00, 0x00, 0x30, 0x61, 0x00, 0x00, 0x30, 0x6f, 0x00, 0x00, 0x4e, 0x16, 0x00, 0x00, 0x75, 0x4c,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x61, 0x00, 0x00, 0x00, 0x62, 0x00, 0x00, 0x00, 0x63,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x7f, 0x00, 0x00, 0x00, 0x21,
	}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "_BE_こ__\nん__に__\nち__は__\n世__界__\n.___a___\nb___c___\n.___.___\n!___\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF32_emoji(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-32", 16)
	input := []byte{
		0x00, 0x00, 0xfe, 0xff, 0x00, 0x01, 0xf4, 0xa9, 0x00, 0x00, 0x30, 0x84, 0x00, 0x01, 0xf3, 0x7a,
		0x00, 0x00, 0x30, 0x6a, 0x00, 0x00, 0x30, 0x69, 0x00, 0x00, 0x30, 0x6e, 0x00, 0x00, 0x7d, 0x75,
		0x00, 0x00, 0x65, 0x87, 0x00, 0x00, 0x5b, 0x57,
	}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "_BE_💩__や__🍺__\nな__ど__の__絵__\n文__字__\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteUTF32_hankaku(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf32", 8)
	input := []byte{
		0x00, 0x00, 0xfe, 0xff, 0x00, 0x00, 0xff, 0x8a, 0x00, 0x00, 0xff, 0x9d, 0x00, 0x00, 0xff, 0x76,
		0x00, 0x00, 0xff, 0x78, 0x00, 0x00, 0xff, 0x76, 0x00, 0x00, 0xff, 0x85, 0x00, 0x00, 0xff, 0x93,
		0x00, 0x00, 0xff, 0x7c, 0x00, 0x00, 0xff, 0x9e,
	}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "_BE_ﾊ___\nﾝ___ｶ___\nｸ___ｶ___\nﾅ___ﾓ___\nｼ___ﾞ___\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteEUCKR(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "euc-kr", 8)
	input := []byte{0xbe, 0xc8, 0xb3, 0xe7, 0xc8, 0xf7, 0xb0, 0xe8, 0xbd, 0xca, 0xbd, 0xc3, 0xbf, 0xc0}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "안녕히계\n십시오\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteEUCCN(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "euc-cn", 8)
	input := []byte{0xc4, 0xe3, 0xba, 0xc3, 0xce, 0xd2, 0xba, 0xc3}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "你好我好\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_WriteBig5(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "big5", 8)
	input := []byte{0xa7, 0x41, 0xa6, 0x6e, 0xa7, 0xda, 0xa6, 0x6e}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "你好我好\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

//nolint:gosmopolitan
func TestPrintable_Write8859_1(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "iso 8859-1", 8)
	input := []byte{0x48, 0xe4, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0xf6, 0x72, 0x6c, 0x64}
	_, err := p.Write(input)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err = p.Close(); err != nil {
		t.Error("close", "err", err)
	}
	expected := "Hä_llo W\nö_rld\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}

func TestPrintable_Close(t *testing.T) {
	buf := &bytes.Buffer{}
	p := NewPrintable(buf, "utf-8", 8)
	_, _ = p.Write([]byte("abc"))
	if err := p.Close(); err != nil {
		t.Errorf("Close error: %v", err)
	}
	expected := "abc\n"
	if buf.String() != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", buf.String(), expected)
	}
}
