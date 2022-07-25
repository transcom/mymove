package testingsuite

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// func (suite *PopTestSuite) SetupTest() {

// }

// func TestPopTestSuite(t *testing.T) {
// 	ps := NewPopTestSuite(CurrentPackage(), WithPerTestTransaction())
// 	suite.Run(t, ps)
// 	ps.TearDown()
// }

// func (suite *PopTestSuite) TestRunWithPreloadedData() {
// 	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

// 	suite.RunWithPreloadedData("Main test records available in subtest", func() {
// 		var foundAddress models.Address
// 		err := suite.DB().Find(&foundAddress, address.ID)
// 		suite.NoError(err)
// 		suite.Equal(address.ID, foundAddress.ID)

// 	})

// 	var address2 models.Address
// 	suite.RunWithPreloadedData("Subtest record creation", func() {
// 		// Create address to search for in the next test
// 		address2 = testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
// 	})

// 	suite.RunWithPreloadedData("Subtest record not found", func() {
// 		// Check that address2 cannot be found
// 		var foundAddress models.Address
// 		err := suite.DB().Find(&foundAddress, address2.ID)
// 		suite.Error(err)
// 		suite.NotEqual(address2.ID, foundAddress.ID)
// 	})

// }

// func (suite *PopTestSuite) TestRun() {
// 	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

// 	suite.Run("Main test records not available in subtest", func() {
// 		var foundAddress models.Address
// 		err := suite.DB().Find(&foundAddress, address.ID)
// 		suite.Error(err)
// 		suite.Contains(err.Error(), "no rows in result set")

// 	})

// 	var address2 models.Address
// 	suite.Run("Subtest record creation", func() {
// 		// Create address to search for in the next test
// 		address2 = testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
// 	})

// 	suite.Run("Subtest record not found", func() {
// 		// Check that address2 cannot be found
// 		var foundAddress models.Address
// 		err := suite.DB().Find(&foundAddress, address2.ID)
// 		suite.Error(err)
// 		suite.NotEqual(address2.ID, foundAddress.ID)
// 	})

// }

// Below functions setup a test suite with additional data loading

type AltPopSuite struct {
	*PopTestSuite
	models.ReServices
}

func (suite *AltPopSuite) SetupTest() {
	fmt.Println("💥Adding ", suite.ReServices[2].Code, suite.ReServices[2].ID, "in SetupTest")
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: suite.ReServices[2],
	})

	fmt.Println("💥Adding ", suite.ReServices[3].Code, suite.ReServices[3].ID, "in SetupTest")
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: suite.ReServices[3],
	})
}

func (suite *AltPopSuite) SetupSuite() {
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
	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

	suite.RunWithPreloadedData("Run a test to check if preloads are there", func() {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.NoError(err)

		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err = suite.DB().Find(&foundReService, reservice.ID)
			fmt.Println(reservice.Code)
			suite.NoError(err, "Reservice %s not found", reservice.Code)

		}

	})
}
func (suite *AltPopSuite) TestRunAltAgain() {
	address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

	suite.RunWithPreloadedData("Run it again", func() {
		var foundAddress models.Address
		err := suite.DB().Find(&foundAddress, address.ID)
		suite.NoError(err)

		var foundReService models.ReService
		for _, reservice := range suite.ReServices {
			err = suite.DB().Find(&foundReService, reservice.ID)
			fmt.Println(reservice.Code)
			suite.NoError(err, "Reservice %s not found", reservice.Code)
		}

	})
}
