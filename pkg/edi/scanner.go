package edi

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/transcom/mymove/pkg/edi/segment"
)

// Scanner scans an io.Reader for EDI segments.  Call Next() to advance to
// the next segment, then call Segment() to get the *Segment found.  If
// Next() returns false, parsing of the segment failed (or there was no more
// data).  Use Err() to determine what happened.
type Scanner struct {
	err      error
	fieldSep rune
	scanner  *bufio.Scanner
	segment  edisegment.Segment
}

// NewScanner creates a Scanner reading from r
func NewScanner(r io.Reader, fieldSep rune) *Scanner {
	return &Scanner{scanner: bufio.NewScanner(r), fieldSep: fieldSep}
}

// Next advances to the next Segment. Returns false if no more data
// is available. Look at Err() to determine why.
func (s *Scanner) Next() bool {
	s.scanner.Scan()
	s.err = s.scanner.Err()
	line := s.scanner.Text()
	if s.err != nil || line == "" {
		return false
	}

	parts := strings.Split(s.scanner.Text(), string(s.fieldSep))
	segID := parts[0]

	var seg edisegment.Segment

	switch segID {
	case "B3":
		seg = &edisegment.B3{}
	case "B3A":
		seg = &edisegment.B3A{}
	case "GE":
		seg = &edisegment.GE{}
	case "GS":
		seg = &edisegment.GS{}
	case "G62":
		seg = &edisegment.G62{}
	case "IEA":
		seg = &edisegment.IEA{}
	case "ISA":
		seg = &edisegment.ISA{}
	case "LX":
		seg = &edisegment.LX{}
	case "L0":
		seg = &edisegment.L0{}
	case "L1":
		seg = &edisegment.L1{}
	case "L5":
		seg = &edisegment.L5{}
	case "L7":
		seg = &edisegment.L7{}
	case "MEA":
		seg = &edisegment.MEA{}
	case "NTE":
		seg = &edisegment.NTE{}
	case "N1":
		seg = &edisegment.N1{}
	case "N4":
		seg = &edisegment.N4{}
	case "N9":
		seg = &edisegment.N9{}
	case "SE":
		seg = &edisegment.SE{}
	case "ST":
		seg = &edisegment.ST{}
	default:
		// TODO - handle better
		panic(fmt.Sprintf("cannot parse segment type %q", segID))
	}

	if err := seg.Parse(parts[1:]); err != nil {
		s.err = err
		return false
	}

	s.segment = seg
	return true
}

// Segment returns the most recent Segment read from the input. Next()
// must be called previously for this value to be valid.
func (s *Scanner) Segment() edisegment.Segment {
	return s.segment
}

// Err returns the most recent error scanning input.
func (s *Scanner) Err() error {
	return s.err
}
