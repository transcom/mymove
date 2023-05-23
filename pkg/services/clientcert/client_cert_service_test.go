package clientcert

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ClientCertServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestUserSuite(t *testing.T) {

	ts := &ClientCertServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
