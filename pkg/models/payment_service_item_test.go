package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentServiceItemValidation() {
	cents := unit.Cents(1000)

	suite.T().Run("test valid PaymentServiceItem", func(t *testing.T) {
		validPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
			MTOServiceItemID: uuid.Must(uuid.NewV4()), //MTO Service Item
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
			PriceCents:       &cents,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentServiceItem, expErrors)
	})

	suite.T().Run("test empty PaymentServiceItem", func(t *testing.T) {
		invalidPaymentServiceItem := models.PaymentServiceItem{}

		expErrors := map[string][]string{
			"payment_request_id": {"PaymentRequestID can not be blank."},
			"mtoservice_item_id": {"MTOServiceItemID can not be blank."},
			"status":             {"Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID, EDI_ERROR]."},
			"requested_at":       {"RequestedAt can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})

	suite.T().Run("test invalid status for PaymentServiceItem", func(t *testing.T) {
		invalidPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()),
			MTOServiceItemID: uuid.Must(uuid.NewV4()), //MTO Service Item
			Status:           "Sleeping",
			RequestedAt:      time.Now(),
			PriceCents:       &cents,
		}
		expErrors := map[string][]string{
			"status": {"Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID, EDI_ERROR]."},
		}
		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})
}

func (suite *ModelSuite) TestPSIBeforeCreate() {
	serviceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{})
	move := serviceItem.MoveTaskOrder
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	suite.T().Run("test with no ID or Reference ID", func(t *testing.T) {
		paymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: paymentRequest.ID,
			MTOServiceItemID: serviceItem.ID,
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}

		err := paymentServiceItem.BeforeCreate(suite.DB())
		suite.NoError(err)
		suite.NotEqual(uuid.Nil, paymentServiceItem.ID)
		suite.NotEmpty(paymentServiceItem.ReferenceID)
	})

	suite.T().Run("test with ID and Reference ID already provided", func(t *testing.T) {
		psiID := uuid.FromStringOrNil("8dce708b-58ab-4adc-a243-ae0c53a44a41")
		referenceID := "1234-5678-8dce708b"
		filledPaymentServiceItem := models.PaymentServiceItem{
			ID:               psiID,
			ReferenceID:      referenceID,
			PaymentRequestID: paymentRequest.ID,
			MTOServiceItemID: serviceItem.ID,
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}

		err := filledPaymentServiceItem.BeforeCreate(suite.DB())
		suite.NoError(err)
		suite.Equal(psiID, filledPaymentServiceItem.ID)
		suite.Equal(referenceID, filledPaymentServiceItem.ReferenceID)
	})

	suite.T().Run("test failure because payment request not found", func(t *testing.T) {
		badPaymentServiceItem := models.PaymentServiceItem{
			PaymentRequestID: uuid.Must(uuid.NewV4()), // new UUID pointing nowhere
			MTOServiceItemID: serviceItem.ID,
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}

		err := badPaymentServiceItem.BeforeCreate(suite.DB())
		suite.Error(err)
	})
}

func (suite *ModelSuite) TestGeneratePSIReferenceID() {
	serviceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{})
	move := serviceItem.MoveTaskOrder
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: move,
	})
	mtoReferenceID := move.ReferenceID
	suite.NotNil(mtoReferenceID)

	baseIDMinusTwoDigits := "caaa0192-c41a-4448-b023-01e15e8fd1"
	paymentServiceItem := models.PaymentServiceItem{
		ID:               uuid.FromStringOrNil(baseIDMinusTwoDigits + "00"),
		PaymentRequestID: paymentRequest.ID,
		MTOServiceItemID: serviceItem.ID,
		Status:           "REQUESTED",
		RequestedAt:      time.Now(),
	}

	// Another PSI with am ID that differs by a digit.
	paymentServiceItem2 := models.PaymentServiceItem{
		ID:               uuid.FromStringOrNil(baseIDMinusTwoDigits + "01"),
		PaymentRequestID: paymentRequest.ID,
		MTOServiceItemID: uuid.Must(uuid.NewV4()),
		Status:           "REQUESTED",
		RequestedAt:      time.Now(),
	}

	suite.T().Run("test normal reference ID generation", func(t *testing.T) {
		referenceID, err := paymentServiceItem.GeneratePSIReferenceID(suite.DB())
		suite.NoError(err)

		psiIDDigits := fmt.Sprintf("%x", paymentServiceItem.ID)
		suite.Equal(*mtoReferenceID+"-"+psiIDDigits[:models.PaymentServiceItemMinReferenceIDSuffixLength], referenceID)

		paymentServiceItem.ReferenceID = referenceID
		suite.MustCreate(suite.DB(), &paymentServiceItem)
	})

	suite.T().Run("test another payment request with ID that differs by a digit", func(t *testing.T) {
		psiIDDigits := fmt.Sprintf("%x", paymentServiceItem2.ID)

		referenceID, err := paymentServiceItem2.GeneratePSIReferenceID(suite.DB())
		suite.NoError(err)
		suite.Equal(*mtoReferenceID+"-"+psiIDDigits[:models.PaymentServiceItemMinReferenceIDSuffixLength+1], referenceID)
	})

	suite.T().Run("test running out of hex digits", func(t *testing.T) {
		// Need to create PSIs up to the max hex digits -- we already have the first one from above.
		start := models.PaymentServiceItemMinReferenceIDSuffixLength + 1
		end := models.PaymentServiceItemMaxReferenceIDLength - len(*mtoReferenceID) - 1
		if end >= len(baseIDMinusTwoDigits)+2 {
			end = len(baseIDMinusTwoDigits)
		}

		for i := start; i <= end; i++ {
			testID := uuid.FromStringOrNil(baseIDMinusTwoDigits + fmt.Sprintf("%02d", i))
			psiIDDigits := fmt.Sprintf("%x", testID)
			longRefID := *mtoReferenceID + "-" + psiIDDigits[:i]
			psiLongReferenceID := models.PaymentServiceItem{
				ID:               testID,
				PaymentRequestID: paymentRequest.ID,
				MTOServiceItemID: serviceItem.ID,
				ReferenceID:      longRefID,
				Status:           "REQUESTED",
				RequestedAt:      time.Now(),
			}
			suite.MustCreate(suite.DB(), &psiLongReferenceID)
		}

		_, err := paymentServiceItem2.GeneratePSIReferenceID(suite.DB())
		suite.Error(err)
		suite.Equal("cannot find unique PSI reference ID", err.Error())
	})
}
