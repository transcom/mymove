package sitentrydateupdate

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UpdateSitEntryDateServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestUpdateSitEntryDateServiceSuite(t *testing.T) {

	hs := &UpdateSitEntryDateServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
