package paymentserviceitem

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaymentServiceItemSuite struct {
	*testingsuite.PopTestSuite
}

func TestPaymentServiceItemServiceSuite(t *testing.T) {
	ts := &PaymentServiceItemSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
