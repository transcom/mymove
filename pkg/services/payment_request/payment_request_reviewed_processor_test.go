package paymentrequest

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const testDateFormat = "060102"

func (suite *PaymentRequestServiceSuite) createPaymentRequest(num int) {

	for i := 0; i < 4; i++ {
		currentTime := time.Now()
		basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				//Value:   currentTime.Format(dateFormat),
				Value: currentTime.Format(testDateFormat),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "4242",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "2424",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip5,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "24245",
			},
		}

		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: mto,
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewed,
				RejectionReason: nil,
			},
		})

		requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)
		scheduledPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 20, 0, 0, 0, 0, time.UTC)
		actualPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 22, 0, 0, 0, 0, time.UTC)

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
				ScheduledPickupDate: &scheduledPickupDate,
				ActualPickupDate:    &actualPickupDate,
			},
		})

		assertions := testdatagen.Assertions{
			Move:           mto,
			MTOShipment:    mtoShipment,
			PaymentRequest: paymentRequest,
		}

		// dlh
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
			assertions,
		)
		// fsc
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeFSC,
			basicPaymentServiceItemParams,
			assertions,
		)
		// ms
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeMS,
			basicPaymentServiceItemParams,
			assertions,
		)
		// cs
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeCS,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dsh
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDSH,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dop
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOP,
			basicPaymentServiceItemParams,
			assertions,
		)
		// ddp
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDP,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dpk
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDPK,
			basicPaymentServiceItemParams,
			assertions,
		)
		// dupk
		_ = testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDUPK,
			basicPaymentServiceItemParams,
			assertions,
		)
	}
}

func (suite *PaymentRequestServiceSuite) TestProcessReviewedPaymentRequest() {

	suite.T().Run("process reviewed payment request succesfully", func(t *testing.T) {

		suite.createPaymentRequest(4)

		_ = testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				ID: uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
			},
		})

		generator := invoice.NewGHCPaymentRequestInvoiceGenerator(suite.DB())

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(suite.DB(), suite.logger, generator)
		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.NoError(err)
	})

	suite.T().Run("process reviewed payment request, failed EDI generator", func(t *testing.T) {

		suite.createPaymentRequest(4)

		// ediinvoice.Invoice858C, error
		ediGenerator := &mocks.GHCPaymentRequestInvoiceGenerator{}
		ediGenerator.
			On("Generate", mock.Anything, mock.Anything).Return(ediinvoice.Invoice858C{}, errors.New("test error"))

		// Process Reviewed Payment Requests
		paymentRequestReviewedProcessor := NewPaymentRequestReviewedProcessor(suite.DB(), suite.logger, ediGenerator)
		err := paymentRequestReviewedProcessor.ProcessReviewedPaymentRequest()
		suite.Contains(err.Error(), "function ProcessReviewedPaymentRequest failed call")
	})
}
