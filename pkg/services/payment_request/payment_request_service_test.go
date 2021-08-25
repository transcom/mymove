package paymentrequest

import (
	"testing"

	"github.com/spf13/afero"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// PaymentRequestServiceSuite is a suite for testing payment requests
type PaymentRequestServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
	fs     *afero.Afero
}

// TestAppContext returns the AppContext for the test suite
func (suite *PaymentRequestServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.AppContextTestHelper.CurrentTestContext(suite.T().Name()), suite.DB())
}

func TestPaymentRequestServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &PaymentRequestServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing,
		fs:                   file,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
