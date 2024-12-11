package testingsuite

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

type SimplePopSuite struct {
	*PopTestSuite
	models.ReServices
}

func (suite *SimplePopSuite) SetupTest() {

}

func TestSimplePopSuite(t *testing.T) {
	sp := &SimplePopSuite{
		PopTestSuite: NewPopTestSuite(CurrentPackage(), WithPerTestTransaction()),
	}

	suite.Run(t, sp)
	sp.TearDown()
}

func (suite *SimplePopSuite) TestRunWithPreloadData() {
	var address models.Address
	suite.PreloadData(func() {
		address = factory.BuildAddress(suite.DB(), nil, nil)
	})

	suite.Run("PreloadData test records available in subtest", func() {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.NoError(err)
		suite.Equal(address.ID, foundAddress.ID)

	})

	var address2 models.Address
	suite.Run("Subtest record creation", func() {
		// Create address to search for in the next test
		address2 = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
	})

	suite.Run("Subtest record not found", func() {
		// Check that address2 cannot be found
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address2.ID)
		suite.Error(err)
		suite.NotEqual(address2.ID, foundAddress.ID)
	})

	suite.Run("non testify subtest", func() {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.NoError(err)
		suite.Equal(address.ID, foundAddress.ID)
	})

}

func (suite *SimplePopSuite) TestMultipleDBPanic() {
	address := factory.BuildAddress(suite.DB(), nil, nil)

	suite.Run("Trying to use db in main and subtest panics", func() {
		panicFunc := func() {
			var foundAddress models.Address
			err := suite.DB().Find(&foundAddress, address.ID) // should panic
			suite.Error(err)
		}
		suite.Panics(panicFunc)
	})
}

func (suite *SimplePopSuite) TestRun() {
	var address2 models.Address
	suite.Run("Subtest record creation", func() {
		// Create address to search for in the next test
		address2 = factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
	})

	suite.Run("Subtest record not found", func() {
		// Check that address2 cannot be found
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address2.ID)
		suite.Error(err)
		suite.NotEqual(address2.ID, foundAddress.ID)
	})

}

type PreloadedPopSuite struct {
	*PopTestSuite
	models.ReServices
}

func (suite *PreloadedPopSuite) SetupTest() {

}

func (suite *PreloadedPopSuite) TearDownSuite() {
	suite.PopTestSuite.TearDown()
}

func TestPreloadedPopSuite(t *testing.T) {

	reservices := []models.ReService{
		{
			Code: models.ReServiceCodeCS,
		},
		{
			Code: models.ReServiceCodeMS,
		},
		{
			Code: models.ReServiceCodeDCRT,
		},
	}

	hs := &PreloadedPopSuite{
		PopTestSuite: NewPopTestSuite(CurrentPackage(), WithPerTestTransaction()),
		ReServices:   reservices,
	}

	suite.Run(t, hs)
}

func (suite *PreloadedPopSuite) TestRunAlt() {

	suite.Run("Run a test to check if preloads are available", func() {
		// Under test:       suite.PreloadData
		// Set up:           This suite has preloaded data in the SetupSuite function
		//                   We only add one new reService in this subtest
		// Expected outcome: The 3 preloaded reService items are found in the database

		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Where("code = $1", reservice.Code).First(&foundReService)
			suite.NoError(err, "Reservice %s not found", reservice.Code)
		}
		// Add a DUCRT ReService, this should not exist outside this subtest
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUCRT)

	})
	suite.Run("Run a test to check that subtests are isolated", func() {
		// Under test:       suite.PreloadData
		// Set up:           This suite has preloaded data in the SetupSuite function
		// Expected outcome: The 3 preloaded reService items are found in the database
		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Where("code = $1", reservice.Code).First(&foundReService)
			suite.NoError(err, "Reservice %s not found", reservice.Code)
		}
	})

}
func (suite *PreloadedPopSuite) TestRunAltAgain() {

	suite.Run("Run a second test to check if preloads are available", func() {
		// Under test:       suite.PreloadData
		// Set up:           This suite has preloaded data in the SetupSuite function
		// Expected outcome: The 3 preloaded reService items are found in the database

		// Reason for a second test is to ensure they don't accidentally get cleaned up between tests
		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Where("code = $1", reservice.Code).First(&foundReService)
			suite.NoError(err, "Reservice %s not found", reservice.Code)

		}
	})
}
