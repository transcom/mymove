package models_test

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentServiceItemValidation() {
	cents := unit.Cents(1000)

	suite.Run("test valid PaymentServiceItem", func() {
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

	suite.Run("test empty PaymentServiceItem", func() {
		invalidPaymentServiceItem := models.PaymentServiceItem{}

		expErrors := map[string][]string{
			"payment_request_id": {"PaymentRequestID can not be blank."},
			"mtoservice_item_id": {"MTOServiceItemID can not be blank."},
			"status":             {"Status is not in the list [REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID, EDI_ERROR]."},
			"requested_at":       {"RequestedAt can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPaymentServiceItem, expErrors)
	})

	suite.Run("test invalid status for PaymentServiceItem", func() {
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

	suite.Run("test with no ID or Reference ID", func() {
		serviceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{})
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: serviceItem.MoveTaskOrder,
		})
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

	suite.Run("test with ID and Reference ID already provided", func() {
		serviceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{})
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: serviceItem.MoveTaskOrder,
		})
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

	suite.Run("test failure because payment request not found", func() {
		serviceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: serviceItem.MoveTaskOrder,
		})
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

	setupTestData := func() (models.PaymentRequest, models.MTOServiceItem) {
		serviceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{})
		move := serviceItem.MoveTaskOrder
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		fmt.Println("paymentRequestID", paymentRequest.ID)
		fmt.Println("moveTaskOrderID", paymentRequest.MoveTaskOrderID)
		fmt.Println("moveReferenceID", *paymentRequest.MoveTaskOrder.ReferenceID)

		return paymentRequest, serviceItem
	}

	suite.Run("test normal reference ID generation", func() {
		// Under test:       GeneratePSIReferenceID returns a reference ID for the PaymentServiceItem it is being called on.
		// Mocked:
		// Set up:           Create a PSI, then generate a reference id.
		// Expected outcome: Generated referenceID matches the expected pattern of reference id + psiIDDigits
		paymentRequest, serviceItem := setupTestData()

		// We want both PSIs to have similar UUIDs, so we force them to be similar
		baseIDMinusTwoDigits := "caaa0192-c41a-4448-b023-01e15e8fd1"
		paymentServiceItem1ID := uuid.FromStringOrNil(baseIDMinusTwoDigits + "00")
		paymentServiceItem2ID := uuid.FromStringOrNil(baseIDMinusTwoDigits + "01")

		// Create first psi
		paymentServiceItem := models.PaymentServiceItem{
			ID:               paymentServiceItem1ID,
			PaymentRequestID: paymentRequest.ID,
			MTOServiceItemID: serviceItem.ID,
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}
		// Generate ID
		referenceID, err := paymentServiceItem.GeneratePSIReferenceID(suite.DB())
		suite.NoError(err)
		fmt.Println("paymentServiceItemID", paymentServiceItem.ID.String())

		// Calculated expected PSI reference id (combo of mtoReferenceID + first N digits of PaymentServiceID)
		mtoReferenceID := serviceItem.MoveTaskOrder.ReferenceID
		psiIDDigits := fmt.Sprintf("%x", paymentServiceItem.ID)
		psiIDDigits = psiIDDigits[:models.PaymentServiceItemMinReferenceIDSuffixLength]
		suite.Equal(*mtoReferenceID+"-"+psiIDDigits, referenceID)

		// Store the first PSI
		paymentServiceItem.ReferenceID = referenceID
		suite.MustCreate(&paymentServiceItem)

		// Now we check that a second PSI with similar ID gets a different reference ID
		paymentServiceItem2 := models.PaymentServiceItem{
			ID:               paymentServiceItem2ID,
			PaymentRequestID: paymentRequest.ID,
			MTOServiceItemID: uuid.Must(uuid.NewV4()),
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}
		// Generate ID
		referenceID, err = paymentServiceItem.GeneratePSIReferenceID(suite.DB())
		suite.NoError(err)

		// Calculated expected PSI reference id (combo of mtoReferenceID + first N PLUS ONE digits of PaymentServiceID)
		mtoReferenceID = serviceItem.MoveTaskOrder.ReferenceID
		psiIDDigits = fmt.Sprintf("%x", paymentServiceItem2.ID)
		psiIDDigits = psiIDDigits[:models.PaymentServiceItemMinReferenceIDSuffixLength+1]
		suite.Equal(*mtoReferenceID+"-"+psiIDDigits, referenceID)

	})

	suite.Run("test running out of hex digits", func() {
		// Under test:       GeneratePSIReferenceID returns a reference ID for the PaymentServiceItem it is being called on.
		// Mocked:
		// Set up:           Create a PSI, then generate a reference id.
		//                   Then create PSIs so as to max out the possible reference ids.
		//                   Then create one more PSI
		// Expected outcome: Error due to running out of possible reference ids.

		paymentRequest, serviceItem := setupTestData()

		// We want both PSIs to have similar UUIDs, so we force them to be similar
		baseIDMinusTwoDigits := "caaa0192-c41a-4448-b023-01e15e8fd1"
		paymentServiceItem1ID := uuid.FromStringOrNil(baseIDMinusTwoDigits + "00")
		// Create first psi
		paymentServiceItem := models.PaymentServiceItem{
			ID:               paymentServiceItem1ID,
			PaymentRequestID: paymentRequest.ID,
			MTOServiceItemID: serviceItem.ID,
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}
		suite.MustCreate(&paymentServiceItem)

		// Need to create PSIs up to the max hex digits -- we already have the first one from above.
		mtoReferenceID := serviceItem.MoveTaskOrder.ReferenceID
		start := models.PaymentServiceItemMinReferenceIDSuffixLength + 1
		end := models.PaymentServiceItemMaxReferenceIDLength - len(*mtoReferenceID)
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
			suite.MustCreate(&psiLongReferenceID)
		}

		// Now we check that a second PSI with similar ID gets a different reference ID
		paymentServiceItem2ID := uuid.FromStringOrNil(baseIDMinusTwoDigits + "01")
		finalPaymentServiceItem := models.PaymentServiceItem{
			ID:               paymentServiceItem2ID,
			PaymentRequestID: paymentRequest.ID,
			MTOServiceItemID: uuid.Must(uuid.NewV4()),
			Status:           "REQUESTED",
			RequestedAt:      time.Now(),
		}
		_, err := finalPaymentServiceItem.GeneratePSIReferenceID(suite.DB())
		suite.Error(err)
		suite.Equal("cannot find unique PSI reference ID", err.Error())
	})
}
