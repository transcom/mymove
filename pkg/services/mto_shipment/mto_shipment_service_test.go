package mtoshipment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOShipmentServiceSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *MTOShipmentServiceSuite) SetupSuite() {
	suite.PreloadData(func() {
		err := factory.DeleteAllotmentsFromDatabase(suite.DB())
		suite.FatalNoError(err)
		factory.SetupDefaultAllotments(suite.DB())
	})
}

func TestMTOShipmentServiceSuite(t *testing.T) {

	ts := &MTOShipmentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
