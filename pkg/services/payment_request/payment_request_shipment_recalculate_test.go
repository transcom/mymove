package paymentrequest

import (
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestRecalculateShipmentPaymentRequestSuccess() {
	// Setup baseline move/shipment/service items data along with needed rate data.
	move, paymentRequestArg := suite.setupRecalculateData()

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("Zip3TransitDistance",
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	paymentRequest, err := creator.CreatePaymentRequest(suite.TestAppContext(), &paymentRequestArg)
	suite.FatalNoError(err)

	// Add a couple of proof of service docs and prime uploads.
	for i := 0; i < 2; i++ {
		proofOfServiceDoc := testdatagen.MakeProofOfServiceDoc(suite.DB(), testdatagen.Assertions{
			ProofOfServiceDoc: models.ProofOfServiceDoc{
				PaymentRequestID: paymentRequest.ID,
			},
		})
		contractor := testdatagen.MakeDefaultContractor(suite.DB())
		testdatagen.MakePrimeUpload(suite.DB(), testdatagen.Assertions{
			PrimeUpload: models.PrimeUpload{
				ProofOfServiceDocID: proofOfServiceDoc.ID,
				ContractorID:        contractor.ID,
			},
		})
		testdatagen.MakePrimeUpload(suite.DB(), testdatagen.Assertions{
			PrimeUpload: models.PrimeUpload{
				ProofOfServiceDocID: proofOfServiceDoc.ID,
				ContractorID:        contractor.ID,
				DeletedAt:           swag.Time(time.Now()),
			},
		})
	}

	// Adjust shipment's original weight to force different pricing on a recalculation.
	mtoShipment := move.MTOShipments[0]
	newWeight := recalculateTestNewOriginalWeight
	mtoShipment.PrimeActualWeight = &newWeight
	suite.MustSave(&mtoShipment)

	// Recalculate the payment request for shipment
	statusUpdater := NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := NewPaymentRequestRecalculator(creator, statusUpdater)
	shipmentRecalculator := NewPaymentRequestShipmentRecalculator(recalculator)

	var newPaymentRequests *models.PaymentRequests
	newPaymentRequests, err = shipmentRecalculator.ShipmentRecalculatePaymentRequest(suite.TestAppContext(), mtoShipment.ID)
	suite.NoError(err, "successfully recalculated shipment's payment request")
	suite.Equal(1, len(*newPaymentRequests))

	// Fetch the old payment request again -- status should have changed and it should also
	// have proof of service docs now.  Need to eager fetch some related data to use in test
	// assertions below.
	var oldPaymentRequest models.PaymentRequest
	err = suite.DB().
		EagerPreload(
			"PaymentServiceItems.MTOServiceItem.ReService",
			"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
			"ProofOfServiceDocs.PrimeUploads",
		).
		Find(&oldPaymentRequest, paymentRequest.ID)
	suite.FatalNoError(err)

	var newPaymentRequest models.PaymentRequest
	err = suite.DB().
		EagerPreload(
			"PaymentServiceItems.MTOServiceItem.ReService",
			"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
			"ProofOfServiceDocs.PrimeUploads",
		).
		Where("recalculation_of_payment_request_id=?", paymentRequest.ID).
		First(&newPaymentRequest)
	suite.FatalNoError(err)

	// Verify some top-level items on the payment requests.
	suite.Equal(oldPaymentRequest.MoveTaskOrderID, newPaymentRequest.MoveTaskOrderID, "Both payment requests should point to same move")
	suite.Len(oldPaymentRequest.PaymentServiceItems, 4)
	suite.Equal(len(oldPaymentRequest.PaymentServiceItems), len(newPaymentRequest.PaymentServiceItems), "Both payment requests should have same number of service items")
	suite.Equal(oldPaymentRequest.Status, models.PaymentRequestStatusDeprecated, "Old payment request status incorrect")
	suite.Equal(newPaymentRequest.Status, models.PaymentRequestStatusPending, "New payment request status incorrect")

	// Make sure the links between payment requests are set up properly.
	suite.Nil(oldPaymentRequest.RecalculationOfPaymentRequestID, "Old payment request should have nil link")
	if suite.NotNil(newPaymentRequest.RecalculationOfPaymentRequestID, "New payment request should not have nil link") {
		suite.Equal(oldPaymentRequest.ID, *newPaymentRequest.RecalculationOfPaymentRequestID, "New payment request should link to the old payment request ID")
	}
}

func (suite *PaymentRequestServiceSuite) TestRecalculateShipmentPaymentRequestErrors() {
	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("Zip3TransitDistance",
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := NewPaymentRequestRecalculator(creator, statusUpdater)
	shipmentRecalculator := NewPaymentRequestShipmentRecalculator(recalculator)

	suite.T().Run("Fail to find shipment ID", func(t *testing.T) {
		bogusShipmentID := uuid.Must(uuid.NewV4())

		_, err := shipmentRecalculator.ShipmentRecalculatePaymentRequest(suite.TestAppContext(), bogusShipmentID)
		suite.NoError(err) // Not find a shipment ID doesn't not produce an error. Simply no payment requests are found
		// and nil is returned.
		//var nilPaymentReqeusts *models.PaymentRequests
		//suite.Equal(&models.PaymentRequests{}, newPaymentRequests)
		//suite.Nil(newPaymentRequests)
	})

	/*
		suite.T().Run("Old payment status has unexpected status", func(t *testing.T) {
			paidPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					Status: models.PaymentRequestStatusPaid,
				},
			})

			paymentRequest, err := creator.CreatePaymentRequest(suite.TestAppContext(), &paymentRequestArg)
			suite.FatalNoError(err)
		    // Update to PAID

			//newPR, err = p.paymentRequestRecalculator.RecalculatePaymentRequest(txnAppCtx, pr)
			mockPlanner := &mocks.PaymentRequestRecalculator{}
			mockPlanner.On("RecalculatePaymentRequest",
				suite.TestAppContext(),
				paidPaymentRequest,
			).Return(nil, services.NewQueryError("PaymentRequest", fmt.Errorf("testing"), fmt.Sprintf("unexpected error while testing payment request ID %s", paidPaymentRequest.ID.String())))

			//err := suite.DB().Load(&paidPaymentRequest, "PaymentServiceItems.MTOServiceItem.MTOShipment",
			//	"PaymentServiceItems.MTOServiceItem")
			//suite.NoError(err)

			var oldPaymentRequest models.PaymentRequest
			err := suite.DB().
				EagerPreload(
					"PaymentServiceItems.MTOServiceItem.MTOShipment",
				).
				Find(&oldPaymentRequest, paidPaymentRequest.ID)
			suite.FatalNoError(err)

			var newPaymentRequests *models.PaymentRequests
			newPaymentRequests, err = shipmentRecalculator.ShipmentRecalculatePaymentRequest(suite.TestAppContext(), *oldPaymentRequest.PaymentServiceItems[0].MTOServiceItem.MTOShipmentID)
			suite.NoError(err)
			if suite.Error(err) {
				suite.IsType(services.ConflictError{}, err)
				suite.Contains(err.Error(), paidPaymentRequest.ID.String())
				suite.Contains(err.Error(), models.PaymentRequestStatusPaid)
			}
			suite.Nil(newPaymentRequests)

		})
	*/

	/*
		suite.T().Run("Can handle error when creating new recalculated payment request", func(t *testing.T) {
			errString := "mock payment request recalculate test error"
			mockRecalculator := &mocks.PaymentRequestRecalculator{}
			mockRecalculator.On("RecalculatePaymentRequest",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PaymentRequest"),
				mock.AnythingOfType("bool"),
			).Return(nil, errors.New(errString))

			shipmentRecalculatorWithMockRecalculate := NewPaymentRequestShipmentRecalculator(mockRecalculator)

			paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

			suite.DB().Load(paymentRequest,"PaymentServiceItems.MTOServiceItem.MTOShipment",
				"PaymentServiceItems.MTOServiceItem")

			err := shipmentRecalculatorWithMockRecalculate.ShipmentRecalculatePaymentRequest(suite.TestAppContext(), *paymentRequest.PaymentServiceItems[0].MTOServiceItem.MTOShipmentID)
			if suite.Error(err) {
				suite.Equal(err.Error(), errString)
			}

		})

	*/
}

/*
func (suite *PaymentRequestServiceSuite) Test_findPendingPaymentRequestsForShipment() {

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("Zip3TransitDistance",
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Mock out payment request recalculate
	mockRecalculator := &mocks.PaymentRequestRecalculator{}
	shipmentRecalculator := NewPaymentRequestShipmentRecalculator(mockRecalculator)

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

	_ = testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusEDIError,
			},
		})

	_ = testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
		})

	suite.T().Run("No available payment request to recalculate", func(t *testing.T) {
		shipmentRecalculator.findPendingPaymentRequestsForShipment
	})

	_ = testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
		})
}
*/
