package testingsuite

import (
	"strings"

	"github.com/stretchr/testify/suite"
)

// LocalTestSuite is a suite for testing
type LocalTestSuite struct {
	suite.Suite
}

// FatalNil ends a test if an object is not nil
func (suite *LocalTestSuite) FatalNil(object interface{}, messages ...string) {
	t := suite.T()
	t.Helper()
	if !suite.Nil(object) {
		if len(messages) > 0 {
			t.Fatal(strings.Join(messages, ","))
		} else {
			t.Fatal()
		}
	}
}

// FatalNoError ends a test if an error is not nil
func (suite *LocalTestSuite) FatalNoError(err error, messages ...string) {
	t := suite.T()
	t.Helper()
	if !suite.NoError(err) {
		if len(messages) > 0 {
			t.Fatalf("%s: %s", strings.Join(messages, ","), err.Error())
		} else {
			t.Fatal(err.Error())
		}
	}
}

// FatalFalse ends a test if a value is not false
func (suite *LocalTestSuite) FatalFalse(b bool, messages ...string) {
	t := suite.T()
	t.Helper()
	if !suite.False(b) {
		if len(messages) > 0 {
			t.Fatalf("%s", strings.Join(messages, ","))
		} else {
			t.Fatal()
		}
	}
}
