package paymentrequest

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"

	"go.uber.org/zap"
)

type PaymentRequestHelperSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *PaymentRequestHelperSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.AppContextTestHelper.CurrentTestContext(suite.T().Name()), suite.DB())
}

func TestPaymentRequestHelperSuite(t *testing.T) {
	ts := &PaymentRequestHelperSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
