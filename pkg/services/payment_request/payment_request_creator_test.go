package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {

	paymentRequest := models.PaymentRequest{
		IsFinal: false,
		Status:  "PENDING",
	}

	creator := NewPaymentRequestCreator(suite.DB())

	// Happy path
	suite.T().Run("Payment request is created successfully", func(t *testing.T) {
		_, verrs, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.NoError(err)
		suite.NoVerrs(verrs)
	})

	// Bad move task order ID
	//suite.T().Run("Given an non-existent move task order id, the create should fail", func(t *testing.T) {
	//	mtoID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
	//	invalidPaymentRequest := models.PaymentRequest{
	//		MoveTaskOrderID: mtoID,
	//		IsFinal:         false,
	//	}
	//	_, _, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
	//	suite.Error(err)
	//})
}
