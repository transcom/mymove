package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: moveTaskOrder.ID,
		IsFinal:         false,
	}

	// Happy path
	suite.T().Run("Payment request is created successfully", func(t *testing.T) {

		creator := NewPaymentRequestCreator(suite.DB())
		_, verrs, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.NoError(err)
		suite.NoVerrs(verrs)
	})
}
