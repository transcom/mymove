package testingsuite

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
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

func (suite *PopTestSuite) TestRunWithPreloadData() {
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

func (suite *PopTestSuite) TestMultipleDBPanic() {
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

func (suite *PopTestSuite) TestRun() {
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

type AltPopSuite struct {
	*PopTestSuite
	models.ReServices
}

func (suite *AltPopSuite) SetupTest() {

}

func (suite *AltPopSuite) SetupSuite() {

	suite.PreloadData(func() {
		// Loads some data into database
		// ReServiceCodeCS
		fmt.Println("💥Adding ", suite.ReServices[0].Code, suite.ReServices[0].ID, "in SetupSuite")
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[0],
		})

		// ReServiceCodeMS
		fmt.Println("💥Adding ", suite.ReServices[1].Code, suite.ReServices[1].ID, "in SetupSuite")
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[1],
		})
		fmt.Println("💥Adding ", suite.ReServices[2].Code, suite.ReServices[2].ID, "in SetupTest")
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[2],
		})

		fmt.Println("💥Adding ", suite.ReServices[3].Code, suite.ReServices[3].ID, "in SetupTest")
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: suite.ReServices[3],
		})
	})

}

func (suite *AltPopSuite) TearDownSuite() {
	suite.PopTestSuite.TearDown()
}

func TestAltPopSuite(t *testing.T) {
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
		{
			ID:   uuid.Must(uuid.NewV4()),
			Code: models.ReServiceCodeDUCRT,
		},
	}

	hs := &AltPopSuite{
		PopTestSuite: NewPopTestSuite(CurrentPackage(), WithPerTestTransaction()),
		ReServices:   reservices,
	}

	suite.Run(t, hs)
}

func (suite *AltPopSuite) TestRunAlt() {

	suite.Run("Run a test to check if preloads are there", func() {

		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Find(&foundReService, reservice.ID)
			fmt.Println(reservice.Code)
			suite.NoError(err, "Reservice %s not found", reservice.Code)

		}
	})
}
func (suite *AltPopSuite) TestRunAltAgain() {

	suite.Run("Run a test to check if preloads are there", func() {

		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err := suite.DB().Find(&foundReService, reservice.ID)
			fmt.Println(reservice.Code)
			suite.NoError(err, "Reservice %s not found", reservice.Code)

		}

	})
}
