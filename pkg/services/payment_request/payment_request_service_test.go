package paymentrequest

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

// PaymentRequestServiceSuite is a suite for testing payment requests
type PaymentRequestServiceSuite struct {
	*testingsuite.PopTestSuite
	fs *afero.Afero
}

func TestPaymentRequestServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &PaymentRequestServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		fs:           file,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
