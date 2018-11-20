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

// NewWriter returns a wrapped csv.Writer with `Comma = '*'`
func NewWriter(w io.Writer) *Writer {
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = '*'
	return &Writer{
		csvWriter,
	}
}

// Write is a wrapper for csv.Write
// Add a single segment to the Writer
func (w *Writer) Write(segment []string) error {
	return w.csv.Write(segment)
}

// WriteAll is a wrapper for csv.WriteAll
// Equivalent of calling Write, Flush, then Error
func (w *Writer) WriteAll(segments [][]string) error {
	return w.csv.WriteAll(segments)
}

// Flush is a wrapper for csv.Flush
// It will flush any segments in the Writer buffer
func (w *Writer) Flush() {
	w.csv.Flush()
}

// Error is a wrapper for csv.Error
// It will retrieve errors from prior calls to the Writer
func (w *Writer) Error() error {
	return w.csv.Error()
}
