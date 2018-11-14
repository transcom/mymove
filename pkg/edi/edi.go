/*
Package edi utilizes Go's Reader/Writer interface
for generating Electronic Data Interchange Files

Right now, this is an incomplete implementation wrapping the csv.Writer,
but with a '*' delimiter.
This helps us:
- Extend the current package if more complex EDI logic is necessary
- Provide a more descriptive name for usage
- Adheres to Go's Writer/Reader interface
*/
package edi

import (
	"encoding/csv"
	"io"
)

// Writer is just a wrapper for csv.Writer
type Writer struct {
	csv *csv.Writer
}

// NewWriter returns a csv.Writer with `Comma = '*'`
func NewWriter(w io.Writer) *Writer {
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = '*'
	return &Writer{
		csvWriter,
	}
}

// Write is a wrapper for csv.Write
func (w *Writer) Write(record []string) error {
	return w.csv.Write(record)
}

// WriteAll is a wrapper for csv.WriteAll
func (w *Writer) WriteAll(records [][]string) error {
	return w.csv.WriteAll(records)
}
