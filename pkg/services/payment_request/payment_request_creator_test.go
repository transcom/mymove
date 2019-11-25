package paymentrequest

import (
	"github.com/gofrs/uuid"
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

	// Bad move task order ID
	suite.T().Run("Given an non-existent move task order id, the create should fail", func(t *testing.T) {
		creator := NewPaymentRequestCreator(suite.DB())
		mtoID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: mtoID,
			IsFinal:         false,
		}
		_, _, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)
	})
}
