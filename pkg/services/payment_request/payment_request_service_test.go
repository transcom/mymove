package paymentrequest

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaymentRequestServiceSuite struct {
	testingsuite.PopTestSuite
}

func TestPaymentRequestServiceSuite(t *testing.T) {

	ts := &PaymentRequestServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
