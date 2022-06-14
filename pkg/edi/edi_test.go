package edi

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type EDISuite struct {
	testingsuite.BaseTestSuite
}

func TestEDISuite(t *testing.T) {
	hs := &EDISuite{}

	suite.Run(t, hs)
}

func (suite *EDISuite) TestNewReader() {
	reader := NewReader(strings.NewReader(""))
	suite.Equal('*', reader.Comma, "Reader.Comma is %c, but should be '*'")
}

func (suite *EDISuite) TestNewWriter() {
	writer := NewWriter(os.Stdout)
	suite.Equal('*', writer.Comma, "Writer.Comma is %c, but should be '*'")
}

func (suite *EDISuite) TestNewScanLine() {
	expected := []string{
		"line 1",
		"line 2",
		"line 3",
		"line 4",
	}

	suite.Run("successfully read lines broken by newline \\n", func() {
		scanner := bufio.NewScanner(strings.NewReader(strings.Join(expected, "\n")))
		scanner.Split(SplitLines)
		idx := 0
		for scanner.Scan() {
			suite.Equal(expected[idx], scanner.Text())
			idx++
		}
		suite.Equal(len(expected), idx, "Processed less lines than expected")
	})

	suite.Run("successfully read lines broken by carriage return and newline \\r\\n", func() {
		scanner := bufio.NewScanner(strings.NewReader(strings.Join(expected, "\r\n")))
		scanner.Split(SplitLines)
		idx := 0
		for scanner.Scan() {
			suite.Equal(expected[idx], scanner.Text())
			idx++
		}
		suite.Equal(len(expected), idx, "Processed less lines than expected")
	})

	suite.Run("successfully read lines broken by only carriage return \\r", func() {
		scanner := bufio.NewScanner(strings.NewReader(strings.Join(expected, "\r")))
		scanner.Split(SplitLines)
		idx := 0
		for scanner.Scan() {
			suite.Equal(expected[idx], scanner.Text())
			idx++
		}
		suite.Equal(len(expected), idx, "Processed less lines than expected")
	})
}
