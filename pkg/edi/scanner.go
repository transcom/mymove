package edi

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Scanner scans an io.Reader for EDI segments.  Call Next() to advance to
// the next segment, then call Segment() to get the *Segment found.  If
// Next() returns false, parsing of the segment failed (or there was no more
// data).  Use Err() to determine what happened.
type Scanner struct {
	err      error
	fieldSep rune
	scanner  *bufio.Scanner
	segment  Segment
}

// Segment represents an EDI segment
type Segment interface {
	String(delimeter string) string
	Parse(parts []string) error
}

// NewScanner creates a Scanner reading from r.
func NewScanner(r io.Reader, fieldSep rune) *Scanner {
	return &Scanner{scanner: bufio.NewScanner(r), fieldSep: fieldSep}
}

// Next advances to the next Segment.  Returns false if no more data
// is available.  Look at Err() to determine why.
func (s *Scanner) Next() bool {
	s.scanner.Scan()
	s.err = s.scanner.Err()
	line := s.scanner.Text()
	if s.err != nil || line == "" {
		return false
	}

	parts := strings.Split(s.scanner.Text(), string(s.fieldSep))
	segID := parts[0]

	var segment Segment

	switch segID {
	case "B3":
		segment = &B3{}
	case "B3A":
		segment = &B3A{}
	case "GE":
		segment = &GE{}
	case "GS":
		segment = &GS{}
	case "G62":
		segment = &G62{}
	case "IEA":
		segment = &IEA{}
	case "ISA":
		segment = &ISA{}
	case "LX":
		segment = &LX{}
	case "L0":
		segment = &L0{}
	case "L1":
		segment = &L1{}
	case "L5":
		segment = &L5{}
	case "L7":
		segment = &L7{}
	case "MEA":
		segment = &MEA{}
	case "NTE":
		segment = &NTE{}
	case "N1":
		segment = &N1{}
	case "N4":
		segment = &N4{}
	case "N9":
		segment = &N9{}
	case "SE":
		segment = &SE{}
	case "ST":
		segment = &ST{}
	default:
		panic(fmt.Sprintf("cannot parse segment type %q", segID))
	}

	if err := segment.Parse(parts[1:]); err != nil {
		s.err = err
		return false
	}

	s.segment = segment
	return true
}

// Segment returns the most recent Segment read from the input.  Next()
// must be called previously for this value to be valid.
func (s *Scanner) Segment() Segment {
	return s.segment
}

// Err returns the most recent error scanning input.
func (s *Scanner) Err() error {
	return s.err
}
