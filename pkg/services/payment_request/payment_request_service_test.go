package paymentrequest

import (
	"testing"

	"github.com/spf13/afero"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// PaymentRequestServiceSuite is a suite for testing payment requests
type PaymentRequestServiceSuite struct {
	testingsuite.PopTestSuite
	logger       Logger
	fs           *afero.Afero
	icnSequencer sequence.Sequencer
}

func (suite *PaymentRequestServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestPaymentRequestServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &PaymentRequestServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(),
		fs:           file,
	}
	ts.icnSequencer = sequence.NewDatabaseSequencer(ts.DB(), ediinvoice.ICNSequenceName)
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
