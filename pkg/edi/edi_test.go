package edi

import (
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
