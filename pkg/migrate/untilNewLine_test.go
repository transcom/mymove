package migrate

import (
	"log"
	"time"
)

func (suite *MigrateSuite) TestUntilNewLine() {

	in := "hello\nworld"
	buf := NewBuffer()

	go func() {
		time.Sleep(time.Millisecond * 1)
		_, err := buf.WriteString(in)
		if err != nil {
			log.Panicf("Error writing string: %s", err)
		}
		buf.Close()
	}()

	wait := 10 * time.Millisecond
	lineNum, out, err := untilNewLine(buf, 0, wait)

	suite.Nil(err)
	suite.Equal(5, lineNum)
	suite.Equal("hello", out)
}

func (suite *MigrateSuite) TestUntilNewLineEOF() {

	in := "hello"
	buf := NewBuffer()

	go func() {
		time.Sleep(time.Millisecond * 1)
		_, err := buf.WriteString(in)
		if err != nil {
			log.Panicf("Error writing string: %s", err)
		}
		buf.Close()
	}()

	wait := 10 * time.Millisecond
	lineNum, out, err := untilNewLine(buf, 0, wait)

	suite.NotNil(err)
	suite.Equal(5, lineNum)
	suite.Equal("hello", out)
}
