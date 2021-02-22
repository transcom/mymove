//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package migrate

import (
	"time"
)

func (suite *MigrateSuite) TestUntilSpace() {

	in := "hello world"
	buf := NewBuffer()

	go func() {
		time.Sleep(time.Millisecond * 1)
		buf.WriteString(in)
		buf.Close()
	}()

	wait := 10 * time.Millisecond
	lineNum, out, err := untilSpace(buf, 0, wait)

	suite.Nil(err)
	suite.Equal(5, lineNum)
	suite.Equal("hello", out)
}

func (suite *MigrateSuite) TestUntilSpaceEOF() {

	in := "hello"
	buf := NewBuffer()

	go func() {
		time.Sleep(time.Millisecond * 1)
		buf.WriteString(in)
		buf.Close()
	}()

	wait := 10 * time.Millisecond
	lineNum, out, err := untilSpace(buf, 0, wait)

	suite.NotNil(err)
	suite.Equal(5, lineNum)
	suite.Equal("hello", out)
}
