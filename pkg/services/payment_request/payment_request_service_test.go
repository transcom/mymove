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
	logger *zap.Logger
	fs     *afero.Afero
}

// TestAppContext returns the AppContext for the test suite
func (suite *PaymentRequestServiceSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func TestPaymentRequestServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &PaymentRequestServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(),
		fs:           file,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
