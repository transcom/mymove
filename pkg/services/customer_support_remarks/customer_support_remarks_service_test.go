package customersupportremarks

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type CustomerSupportRemarksSuite struct {
	*testingsuite.PopTestSuite
}

func TestOfficeMoveRemarksServiceSuite(t *testing.T) {
	ts := CustomerSupportRemarksSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, &ts)
	ts.PopTestSuite.TearDown()
}
