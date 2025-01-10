package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type OrderServiceSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *OrderServiceSuite) SetupSuite() {
	suite.PreloadData(func() {
		err := factory.DeleteAllotmentsFromDatabase(suite.DB())
		suite.FatalNoError(err)
		factory.SetupDefaultAllotments(suite.DB())
	})
}

func TestOrderServiceSuite(t *testing.T) {
	ts := &OrderServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
