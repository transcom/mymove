package postalcode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ValidatePostalCodeTestSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ValidatePostalCodeTestSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestValidatePostalCodeTestSuite(t *testing.T) {
	ts := &ValidatePostalCodeTestSuite{
		testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, ts)
}

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode_ValidPostalCode() {
	postalCodeType := services.PostalCodeType("Destination")
	postalCode := "30813"

	validatePostalCode := NewPostalCodeValidator(suite.DB())
	valid, _ := validatePostalCode.ValidatePostalCode(postalCode, postalCodeType)

	suite.True(valid)
}

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode_InvalidPostalCode() {
	postalCodeType := services.PostalCodeType("Destination")
	postalCode := "00000"

	validatePostalCode := NewPostalCodeValidator(suite.DB())
	valid, _ := validatePostalCode.ValidatePostalCode(postalCode, postalCodeType)

	suite.False(valid)
}
