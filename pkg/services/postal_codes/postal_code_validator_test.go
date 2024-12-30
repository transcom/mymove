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

	suite.Run("Postal code is not in zip3_distances table", func() {
		testPostalCode := "30183"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "not found in zip3_distances")
	})

	suite.Run("Contract year cannot be found", func() {
		testPostalCode := "30183"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		testdatagen.MakeZip3Distance(suite.DB(), testdatagen.Assertions{
			Zip3Distance: models.Zip3Distance{
				FromZip3: testPostalCode[:3],
				ToZip3:   "993",
			},
		})

		suite.buildContractYear(testdatagen.GHCTestYear - 1)

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "could not find contract year")
	})

	suite.Run("Postal code is not in re_zip3s table", func() {
		testPostalCode := "30183"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		testdatagen.MakeZip3Distance(suite.DB(), testdatagen.Assertions{
			Zip3Distance: models.Zip3Distance{
				FromZip3: testPostalCode[:3],
				ToZip3:   "993",
			},
		})

		suite.buildContractYear(testdatagen.GHCTestYear)

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "not found in re_zip3s")
	})

	suite.Run("Postal code is not in re_zip5_rate_areas table", func() {
		testPostalCode := "32102"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		testdatagen.MakeZip3Distance(suite.DB(), testdatagen.Assertions{
			Zip3Distance: models.Zip3Distance{
				FromZip3: testPostalCode[:3],
				ToZip3:   "993",
			},
		})

		reContractYear := suite.buildContractYear(testdatagen.GHCTestYear)
		serviceArea := testdatagen.MakeDefaultReDomesticServiceArea(suite.DB())
		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Zip3:                 testPostalCode[:3],
				Contract:             reContractYear.Contract,
				DomesticServiceArea:  serviceArea,
				HasMultipleRateAreas: true,
			},
		})

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.False(valid)
		suite.Error(err)
		suite.IsType(&apperror.UnsupportedPostalCodeError{}, err)
		suite.Contains(err.Error(), "not found in re_zip5_rate_areas")
	})

	suite.Run("Valid postal code for zip3 with single rate area", func() {
		testPostalCode := "30813"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		testdatagen.MakeZip3Distance(suite.DB(), testdatagen.Assertions{
			Zip3Distance: models.Zip3Distance{
				FromZip3: testPostalCode[:3],
				ToZip3:   "993",
			},
		})

		reContractYear := suite.buildContractYear(testdatagen.GHCTestYear)
		serviceArea := testdatagen.MakeDefaultReDomesticServiceArea(suite.DB())
		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Zip3:                testPostalCode[:3],
				Contract:            reContractYear.Contract,
				DomesticServiceArea: serviceArea,
			},
		})

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.True(valid)
		suite.NoError(err)
	})

	suite.Run("Valid postal code for zip3 with multiple rate areas", func() {
		testPostalCode := "32102"
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), testPostalCode, "CNNQ")

		testdatagen.MakeZip3Distance(suite.DB(), testdatagen.Assertions{
			Zip3Distance: models.Zip3Distance{
				FromZip3: testPostalCode[:3],
				ToZip3:   "993",
			},
		})

		reContractYear := suite.buildContractYear(testdatagen.GHCTestYear)
		serviceArea := testdatagen.MakeDefaultReDomesticServiceArea(suite.DB())
		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Zip3:                 testPostalCode[:3],
				Contract:             reContractYear.Contract,
				DomesticServiceArea:  serviceArea,
				HasMultipleRateAreas: true,
			},
		})

		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReContractYear: reContractYear,
		})
		testdatagen.MakeReZip5RateArea(suite.DB(), testdatagen.Assertions{
			ReZip5RateArea: models.ReZip5RateArea{
				Zip5: testPostalCode,
			},
			ReContract: reContractYear.Contract,
			ReRateArea: rateArea,
		})

		valid, err := postalCodeValidator.ValidatePostalCode(suite.AppContextForTest(), testPostalCode)

		suite.True(valid)
		suite.NoError(err)
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
