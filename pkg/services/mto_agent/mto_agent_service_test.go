package mtoagent

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOAgentServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestMTOAgentServiceSuite(t *testing.T) {
	ts := &MTOAgentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
