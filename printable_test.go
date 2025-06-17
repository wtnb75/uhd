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
	expected := "こんにち\nは世界.a\nbc..!.\n"
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
