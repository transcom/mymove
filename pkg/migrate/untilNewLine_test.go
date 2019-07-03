package migrate

import (
	"time"
)

func (suite *MigrateSuite) TestUntilNewLine() {

	in := "hello\nworld"
	buf := NewBuffer()

	go func() {
		time.Sleep(time.Second * 1)
		buf.WriteString(in)
		buf.Close()
	}()

	wait := 10 * time.Second
	lineNum, out, err := untilNewLine(buf, 0, wait)

	suite.Nil(err)
	suite.Equal(5, lineNum)
	suite.Equal("hello", out)
}

func (suite *MigrateSuite) TestUntilNewLineEOF() {

	in := "hello"
	buf := NewBuffer()

	go func() {
		time.Sleep(time.Second * 1)
		buf.WriteString(in)
		buf.Close()
	}()

	wait := 10 * time.Second
	lineNum, out, err := untilNewLine(buf, 0, wait)

	suite.NotNil(err)
	suite.Equal(5, lineNum)
	suite.Equal("hello", out)
}
