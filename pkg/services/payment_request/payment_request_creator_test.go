package paymentrequest

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})
	moveTaskOrderID := moveTaskOrder.ID
	serviceItemID, _ := uuid.NewV4()
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	rejectionReason := "Missing documentation"

	paymentRequest := models.PaymentRequest{
		ID:              paymentRequestID,
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrderID,
		ServiceItemIDs:  []uuid.UUID{serviceItemID},
		RejectionReason: rejectionReason,
	}

	// Happy path
	suite.T().Run("Payment request is created successfully", func(t *testing.T) {

		creator := NewPaymentRequestCreator(suite.DB())
		_, verrs, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.NoError(err)
		suite.NoVerrs(verrs)
	})
}