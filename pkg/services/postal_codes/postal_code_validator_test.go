package postalcode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/services"
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

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode_ValidPostalCode() {
	postalCodeType := services.PostalCodeType("Destination")
	postalCode := "30813"

	validatePostalCode := NewPostalCodeValidator()
	valid, _ := validatePostalCode.ValidatePostalCode(suite.AppContextForTest(), postalCode, postalCodeType)

	suite.True(valid)
}

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode_InvalidPostalCode() {
	postalCodeType := services.PostalCodeType("Destination")
	postalCode := "00000"

	validatePostalCode := NewPostalCodeValidator()
	valid, _ := validatePostalCode.ValidatePostalCode(suite.AppContextForTest(), postalCode, postalCodeType)

	suite.False(valid)
}
