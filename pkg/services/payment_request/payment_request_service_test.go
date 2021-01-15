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
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return value in the file is used for test database teardown
	//RA: Given the database is being reset for unit test use, there are no unexpected states and conditions to account for
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	suite.DB().TruncateAll() // nolint:errcheck
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
