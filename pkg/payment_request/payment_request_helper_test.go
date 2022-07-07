package paymentrequest

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaymentRequestHelperSuite struct {
	*testingsuite.PopTestSuite
}

func TestPaymentRequestHelperSuite(t *testing.T) {
	ts := &PaymentRequestHelperSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
