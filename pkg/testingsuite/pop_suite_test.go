package testingsuite

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PopTestSuite) SetupTest() {

}

func TestPopTestSuite(t *testing.T) {
	ps := NewPopTestSuite(CurrentPackage(), WithPerTestTransaction())
	suite.Run(t, ps)
	ps.TearDown()
}

func (suite *PopTestSuite) TestRunWithPreloadedData() {
	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

	suite.RunWithPreloadedData("Main test records available in subtest", func() {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.NoError(err)
		suite.Equal(address.ID, foundAddress.ID)

	})

	var address2 models.Address
	suite.RunWithPreloadedData("Subtest record creation", func() {
		// Create address to search for in the next test
		address2 = testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
	})

	suite.RunWithPreloadedData("Subtest record not found", func() {
		// Check that address2 cannot be found
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address2.ID)
		suite.Error(err)
		suite.NotEqual(address2.ID, foundAddress.ID)
	})

}

func (suite *PopTestSuite) TestRun() {
	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

	suite.Run("Main test records not available in subtest", func() {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.Error(err)
		suite.Contains(err.Error(), "no rows in result set")

	})

	var address2 models.Address
	suite.Run("Subtest record creation", func() {
		// Create address to search for in the next test
		address2 = testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
	})

	suite.Run("Subtest record not found", func() {
		// Check that address2 cannot be found
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address2.ID)
		suite.Error(err)
		suite.NotEqual(address2.ID, foundAddress.ID)
	})

}
