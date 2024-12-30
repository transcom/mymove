package postalcode

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ValidatePostalCodeTestSuite struct {
	*testingsuite.PopTestSuite
}

func TestValidatePostalCodeTestSuite(t *testing.T) {
	ts := &ValidatePostalCodeTestSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode() {
	mockClock := clock.NewMock()
	mockClock.Set(time.Date(testdatagen.GHCTestYear, 6, 1, 0, 0, 0, 0, time.UTC))
	postalCodeValidator := NewPostalCodeValidator(mockClock)

	suite.Run("Postal code should be at least 5 characters", func() {
		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), "123")

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "less than 5 characters")
	})

	suite.Run("Postal code should only contain digits", func() {
		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), "1234x")

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "should only contain digits")
	})

	suite.Run("Postal code is not in postal_code_to_gblocs table", func() {
		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), "30907")

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "not found in postal_code_to_gblocs")
	})

	suite.Run("Contract year cannot be found", func() {
		testPostalCode := "30183"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		suite.buildContractYear(testdatagen.GHCTestYear - 1)

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "could not find contract year")
	})
}

func (suite *ValidatePostalCodeTestSuite) buildContractYear(testYear int) models.ReContractYear {
	reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
	reContractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			Contract:  reContract,
			StartDate: time.Date(testYear, time.January, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(testYear, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
	})

	return reContractYear
}
