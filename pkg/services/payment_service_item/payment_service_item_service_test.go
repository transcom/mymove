package paymentserviceitem

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaymentServiceItemSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *PaymentServiceItemSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func TestPaymentServiceItemServiceSuite(t *testing.T) {
	ts := &PaymentServiceItemSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
