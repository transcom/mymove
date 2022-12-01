package signedcertification

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SignedCertificationSuite struct {
	*testingsuite.PopTestSuite
}

func TestSignedCertificationServiceSuite(t *testing.T) {
	ts := &SignedCertificationSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
