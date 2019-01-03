/*
Package edi utilizes Go's csv.Reader/csv.Writer
for generating Electronic Data Interchange Files

This helps us:
- Extend the current package if more complex EDI logic is necessary
- Provide a more descriptive name for usage
- Adhere to patterns of Go's stdlib csv.Writer/csv.Reader
*/
package edi

import (
	"encoding/csv"
	"io"
)

// Writer is just a wrapper for csv.Writer
type Writer struct {
	*csv.Writer
}

// NewWriter returns a wrapped csv.Writer with `Comma = '*'`
func NewWriter(w io.Writer) *Writer {
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = '*'
	return &Writer{
		csvWriter,
	}
}

// Reader is just a wrapper for csv.Reader
type Reader struct {
	*csv.Reader
}

// NewReader returns a wrapped csv.Reader with `Comma = '*'`
func NewReader(r io.Reader) *Reader {
	csvReader := csv.NewReader(r)
	csvReader.Comma = '*'
	return &Reader{
		csvReader,
	}
}
