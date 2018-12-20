package edi

import (
	"os"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	reader := NewReader(strings.NewReader(""))
	if reader.Comma != '*' {
		t.Errorf("Reader.Comma is %c, but should be ','", reader.Comma)
	}
}

func TestNewWriter(t *testing.T) {
	writer := NewWriter(os.Stdout)
	if writer.Comma != '*' {
		t.Errorf("Writer.Comma is %x, but should be '*'", writer.Comma)
	}
}
