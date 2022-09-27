package testingsuite

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
		address = testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
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
		address2 = testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
	})

	suite.Run("Subtest record not found", func() {
		// Check that address2 cannot be found
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address2.ID)
		suite.Error(err)
		suite.NotEqual(address2.ID, foundAddress.ID)
	})

	suite.T().Run("non testify subtest", func(t *testing.T) {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.NoError(err)
		suite.Equal(address.ID, foundAddress.ID)
	})

}

func (suite *SimplePopSuite) TestMultipleDBPanic() {
	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

	suite.Run("Trying to use db in main and subtest panics", func() {
		defer func() {
			if r := recover(); r == nil {
				suite.FailNow("Did not panic")
			}
			// manually clean up after recovering from panic in
			// this test.
			for k := range suite.txnTestDb {
				popConn := suite.txnTestDb[k]
				suite.NoError(popConn.Close())
				delete(suite.txnTestDb, k)
			}
		}()
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.Error(err)
	})
}

func (suite *SimplePopSuite) TestRun() {
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

type PreloadedPopSuite struct {
	*PopTestSuite
	models.ReServices
}

func (suite *PreloadedPopSuite) SetupTest() {

}

func (suite *PreloadedPopSuite) SetupSuite() {

	suite.PreloadData(func() {
		// Loads some data into database
		// ReServiceCodeCS
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[0],
		})

		// ReServiceCodeMS
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[1],
		})

		// ReServiceCodeDCRT
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[2],
		})
	})

}

func (suite *PreloadedPopSuite) TearDownSuite() {
	suite.PopTestSuite.TearDown()
}

func TestPreloadedPopSuite(t *testing.T) {
	reservices := []models.ReService{
		{
			ID:   uuid.Must(uuid.NewV4()),
			Code: models.ReServiceCodeCS,
		},
		{
			ID:   uuid.Must(uuid.NewV4()),
			Code: models.ReServiceCodeMS,
		},
		{
			ID:   uuid.Must(uuid.NewV4()),
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
			err := suite.DB().Find(&foundReService, reservice.ID)
			suite.NoError(err, "Reservice %s not found", reservice.Code)
		}
		// Add a DUCRT ReService, this should not exist outside this subtest
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDUCRT,
			},
		})

	})
	suite.Run("Run a test to check that subtests are isolated", func() {
		// Under test:       suite.PreloadData
		// Set up:           This suite has preloaded data in the SetupSuite function
		// Expected outcome: The 3 preloaded reService items are found in the database
		//                   The one new reService added in the subtest above is NOT found
		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Find(&foundReService, reservice.ID)
			suite.NoError(err, "Reservice %s not found", reservice.Code)
		}

		var foundReServices []models.ReService
		err := suite.DB().Where("code = ?", models.ReServiceCodeDUCRT).All(&foundReServices)
		suite.NoError(err)
		suite.Len(foundReServices, 0)
	})

}
func (suite *PreloadedPopSuite) TestRunAltAgain() {

	suite.Run("Run a second test to check if preloads are available", func() {
		// Under test:       suite.PreloadData
		// Set up:           This suite has preloaded data in the SetupSuite function
		// Expected outcome: The 3 preloaded reService items are found in the database
		//                   The one new reService added in the other test in this suite is NOT found

		// Reason for a second test is to ensure they don't accidentally get cleaned up between tests
		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Find(&foundReService, reservice.ID)
			suite.NoError(err, "Reservice %s not found", reservice.Code)

		}

		var foundReServices []models.ReService
		err := suite.DB().Where("code = ?", models.ReServiceCodeDUCRT).All(&foundReServices)
		suite.NoError(err)
		suite.Len(foundReServices, 0)
	})
}
