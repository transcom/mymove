package postalcode

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ValidatePostalCodeTestSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *ValidatePostalCodeTestSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger)
}

func TestValidatePostalCodeTestSuite(t *testing.T) {
	ts := &ValidatePostalCodeTestSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode_ValidPostalCode() {
	postalCodeType := services.PostalCodeType("Destination")
	postalCode := "30813"

	validatePostalCode := NewPostalCodeValidator()
	valid, _ := validatePostalCode.ValidatePostalCode(suite.TestAppContext(), postalCode, postalCodeType)

	suite.True(valid)
}

func (suite *ValidatePostalCodeTestSuite) TestValidatePostalCode_InvalidPostalCode() {
	postalCodeType := services.PostalCodeType("Destination")
	postalCode := "00000"

	validatePostalCode := NewPostalCodeValidator()
	valid, _ := validatePostalCode.ValidatePostalCode(suite.TestAppContext(), postalCode, postalCodeType)

	suite.False(valid)
}
