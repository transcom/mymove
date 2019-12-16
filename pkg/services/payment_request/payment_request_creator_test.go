package paymentrequest

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestCreatePaymentRequest() {
	// Create some records we'll need to link to
	moveTaskOrder := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: moveTaskOrder,
		ReService: models.ReService{
			Code: "DLH",
		},
	})
	mtoServiceItem2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: moveTaskOrder,
		ReService: models.ReService{
			Code: "DOP",
		},
	})
	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         "weight",
			Description: "actual weight",
			Type:        models.ServiceItemParamTypeInteger,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})
	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         "pickup",
			Description: "requested pickup date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	creator := NewPaymentRequestCreator(suite.DB())

	// Happy path
	suite.T().Run("Payment request is created successfully", func(t *testing.T) {
		paymentRequest := models.PaymentRequest{
			MoveTaskOrderID: moveTaskOrder.ID,
			IsFinal:         false,
			PaymentServiceItems: models.PaymentServiceItems{
				{
					ServiceItemID: mtoServiceItem1.ID,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "weight",
							Value:       "3254",
						},
						{
							IncomingKey: "pickup",
							Value:       "2019-12-16",
						},
					},
				},
				{
					ServiceItemID: mtoServiceItem2.ID,
					PaymentServiceItemParams: models.PaymentServiceItemParams{
						{
							IncomingKey: "weight",
							Value:       "7722",
						},
					},
				},
			},
		}

		_, verrs, err := creator.CreatePaymentRequest(&paymentRequest)
		suite.NoError(err)
		suite.NoVerrs(verrs)
	})

	// Bad move task order ID
	suite.T().Run("Given a non-existent move task order id, the create should fail", func(t *testing.T) {
		mtoID, _ := uuid.FromString("0aee14dd-b5ea-441a-89ad-db4439fa4ea2")
		invalidPaymentRequest := models.PaymentRequest{
			MoveTaskOrderID: mtoID,
			IsFinal:         false,
		}
		_, verrs, err := creator.CreatePaymentRequest(&invalidPaymentRequest)
		suite.Error(err)
		suite.NoVerrs(verrs)
	})
}
