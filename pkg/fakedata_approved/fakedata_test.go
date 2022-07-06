package fakedata

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type FakeDataSuite struct {
	*testingsuite.PopTestSuite
}

func TestFakeDataSuite(t *testing.T) {
	hs := &FakeDataSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
