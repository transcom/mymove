package paymentrequest

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestRecalculateShipmentPaymentRequestSuccess() {
	// Setup baseline move/shipment/service items data along with needed rate data.
	move, paymentRequestArg := suite.setupRecalculateData1()
	_, paymentRequestArg2 := suite.setupRecalculateData2(move, paymentRequestArg.PaymentServiceItems[0].MTOServiceItem.MTOShipment)

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())

	// Payment Request 1
	paymentRequest1, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequestArg)
	suite.FatalNoError(err)

	// Add a couple of proof of service docs and prime uploads.
	for i := 0; i < 2; i++ {
		proofOfServiceDoc := testdatagen.MakeProofOfServiceDoc(suite.DB(), testdatagen.Assertions{
			ProofOfServiceDoc: models.ProofOfServiceDoc{
				PaymentRequestID: paymentRequest1.ID,
			},
		})
		contractor := testdatagen.MakeDefaultContractor(suite.DB())
		testdatagen.MakePrimeUpload(suite.DB(), testdatagen.Assertions{
			PrimeUpload: models.PrimeUpload{
				ProofOfServiceDocID: proofOfServiceDoc.ID,
				ContractorID:        contractor.ID,
			},
		})
	}

	// Adjust shipment's original weight to force different pricing on a recalculation.
	mtoShipment := move.MTOShipments[0]
	newWeight := recalculateTestNewOriginalWeight
	mtoShipment.PrimeActualWeight = &newWeight
	suite.MustSave(&mtoShipment)

	// Create additional payment request for same shipment
	// Payment Request 2
	var paymentRequest2 *models.PaymentRequest
	paymentRequest2, err = creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequestArg2)
	suite.FatalNoError(err)

	// Add a couple of proof of service docs and prime uploads.
	for i := 0; i < 2; i++ {
		proofOfServiceDoc := testdatagen.MakeProofOfServiceDoc(suite.DB(), testdatagen.Assertions{
			ProofOfServiceDoc: models.ProofOfServiceDoc{
				PaymentRequestID: paymentRequest2.ID,
			},
		})
		contractor := testdatagen.MakeDefaultContractor(suite.DB())
		testdatagen.MakePrimeUpload(suite.DB(), testdatagen.Assertions{
			PrimeUpload: models.PrimeUpload{
				ProofOfServiceDocID: proofOfServiceDoc.ID,
				ContractorID:        contractor.ID,
			},
		})
	}

	// Recalculate the payment request for shipment
	statusUpdater := NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := NewPaymentRequestRecalculator(creator, statusUpdater)
	shipmentRecalculator := NewPaymentRequestShipmentRecalculator(recalculator)

	var newPaymentRequests *models.PaymentRequests
	newPaymentRequests, err = shipmentRecalculator.ShipmentRecalculatePaymentRequest(suite.AppContextForTest(), mtoShipment.ID)
	suite.NoError(err, "successfully recalculated shipment's payment request")
	suite.Equal(2, len(*newPaymentRequests))

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
		Find(&oldPaymentRequest, paymentRequest1.ID)
	suite.FatalNoError(err)

	var newPaymentRequest models.PaymentRequest
	err = suite.DB().
		EagerPreload(
			"PaymentServiceItems.MTOServiceItem.ReService",
			"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
			"ProofOfServiceDocs.PrimeUploads",
		).
		Where("recalculation_of_payment_request_id=?", paymentRequest1.ID).
		First(&newPaymentRequest)
	suite.FatalNoError(err)

	// Verify some top-level items on the payment requests.
	suite.Equal(oldPaymentRequest.MoveTaskOrderID, newPaymentRequest.MoveTaskOrderID, "Both payment requests should point to same move")
	suite.Len(oldPaymentRequest.PaymentServiceItems, 5)
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
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := NewPaymentRequestRecalculator(creator, statusUpdater)
	shipmentRecalculator := NewPaymentRequestShipmentRecalculator(recalculator)

	// Setup baseline move/shipment/service items data along with needed rate data.
	_ /*move*/, paymentRequestArg := suite.setupRecalculateData1()

	paidPaymentRequest, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequestArg)
	suite.FatalNoError(err)

	suite.T().Run("Fail to find shipment ID", func(t *testing.T) {
		bogusShipmentID := uuid.Must(uuid.NewV4())

		var returnPaymentRequests *models.PaymentRequests
		returnPaymentRequests, err = shipmentRecalculator.ShipmentRecalculatePaymentRequest(suite.AppContextForTest(), bogusShipmentID)
		suite.NoError(err) // Not finding a shipment ID doesn't produce an error. Simply no payment requests are found
		suite.Nil((*models.PaymentRequests)(nil), returnPaymentRequests)
	})

	suite.T().Run("Old payment status has unexpected status", func(t *testing.T) {
		paidPaymentRequest.Status = models.PaymentRequestStatusPaid
		suite.MustSave(paidPaymentRequest)
		// Update to PAID

		mockPlanner := &mocks.PaymentRequestRecalculator{}
		mockPlanner.On("RecalculatePaymentRequest",
			mock.AnythingOfType("*appcontext.appContext"),
			paidPaymentRequest,
		).Return(nil, apperror.NewQueryError("PaymentRequest", fmt.Errorf("testing"), fmt.Sprintf("unexpected error while testing payment request ID %s", paidPaymentRequest.ID.String())))

		var oldPaymentRequest models.PaymentRequest
		err = suite.DB().
			EagerPreload(
				"PaymentServiceItems.MTOServiceItem.MTOShipment",
			).
			Find(&oldPaymentRequest, paidPaymentRequest.ID)
		suite.FatalNoError(err)

		var newPaymentRequests *models.PaymentRequests
		newPaymentRequests, err = shipmentRecalculator.ShipmentRecalculatePaymentRequest(suite.AppContextForTest(), *oldPaymentRequest.PaymentServiceItems[0].MTOServiceItem.MTOShipmentID)
		suite.NoError(err)
		suite.Nil((*models.PaymentRequests)(nil), newPaymentRequests)
	})

	suite.T().Run("Can handle error when creating new recalculated payment request", func(t *testing.T) {
		errString := "mock payment request recalculate test error"
		mockRecalculator := &mocks.PaymentRequestRecalculator{}
		mockRecalculator.On("RecalculatePaymentRequest",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, errors.New(errString))

		shipmentRecalculatorWithMockRecalculate := NewPaymentRequestShipmentRecalculator(mockRecalculator)

		suite.FatalNoError(err)
		pendingPaymentRequest := paidPaymentRequest
		pendingPaymentRequest.Status = models.PaymentRequestStatusPending
		suite.MustSave(pendingPaymentRequest)
		// Update to PENDING

		err = suite.DB().Load(pendingPaymentRequest, "PaymentServiceItems.MTOServiceItem.MTOShipment",
			"PaymentServiceItems.MTOServiceItem")
		suite.NoError(err)

		returnPaymentRequests, err := shipmentRecalculatorWithMockRecalculate.ShipmentRecalculatePaymentRequest(suite.AppContextForTest(), *pendingPaymentRequest.PaymentServiceItems[0].MTOServiceItem.MTOShipmentID)
		if suite.Error(err) {
			suite.Equal(err.Error(), errString)
		}
		suite.Nil((*models.PaymentRequests)(nil), returnPaymentRequests)
	})
}
