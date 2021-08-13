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
	"bytes"
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

// dropCR drops a terminal \r from the data.
// pulled from bufio library
//
// See https://cs.opensource.google/go/go/+/refs/tags/go1.16.7:src/bufio/scan.go;l=336
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// SplitLines is a split function for BufIO library
// it borrows from the default function but add support for files with only
// carriage returns as the line delimiter. Sometimes EDI files have only carriage returns
//
// See https://cs.opensource.google/go/go/+/refs/tags/go1.16.7:src/bufio/scan.go;l=350
func SplitLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	} else if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full carriage return terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	//If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
