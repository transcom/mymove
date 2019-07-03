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
