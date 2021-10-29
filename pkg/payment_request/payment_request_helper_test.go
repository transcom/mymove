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
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *PaymentRequestHelperSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func TestPaymentRequestHelperSuite(t *testing.T) {
	ts := &PaymentRequestHelperSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
