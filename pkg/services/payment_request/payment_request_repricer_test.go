package paymentrequest

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	repriceTestPickupZip            = "30907"
	repriceTestDestinationZip       = "78234"
	repriceTestMSFee                = unit.Cents(25513)
	repriceTestCSFee                = unit.Cents(22399)
	repriceTestDLHPrice             = unit.Millicents(6000)
	repriceTestFSCPrice             = unit.Millicents(277600)
	repriceTestEstimatedWeight      = unit.Pound(3500)
	repriceTestOriginalWeight       = unit.Pound(3652)
	repriceTestNewOriginalWeight    = unit.Pound(3412)
	repriceTestEscalationCompounded = 1.04071
	repriceTestZip3Distance         = 1234
)

func (suite *PaymentRequestServiceSuite) TestRepricePaymentRequestSuccess() {
	// Setup baseline move/shipment/service items data along with needed rate data.
	move, paymentRequestArg := suite.setupRepriceData()

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("Zip3TransitDistance",
		repriceTestPickupZip,
		repriceTestDestinationZip,
	).Return(repriceTestZip3Distance, nil)

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

	// Adjust shipment's original weight to force different pricing on a reprice.
	mtoShipment := move.MTOShipments[0]
	newWeight := repriceTestNewOriginalWeight
	mtoShipment.PrimeActualWeight = &newWeight
	suite.MustSave(&mtoShipment)

	// Reprice the payment request created above.
	repricer := NewPaymentRequestRepricer(creator)
	newPaymentRequest, err := repricer.RepricePaymentRequest(suite.TestAppContext(), paymentRequest.ID)
	suite.FatalNoError(err)

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

	// Verify some top-level items on the payment requests.
	suite.Equal(oldPaymentRequest.MoveTaskOrderID, newPaymentRequest.MoveTaskOrderID, "Both payment requests should point to same move")
	suite.Equal(len(oldPaymentRequest.PaymentServiceItems), len(newPaymentRequest.PaymentServiceItems), "Both payment requests should have same number of service items")
	suite.Equal(oldPaymentRequest.Status, models.PaymentRequestStatusReviewedAllRejected, "Old payment request status incorrect")
	suite.Equal(newPaymentRequest.Status, models.PaymentRequestStatusPending, "New payment request status incorrect")

	// Verify that the IDs of the MTO service items remain the same across both payment requests.
	oldMTOServiceItemIDs := make(map[uuid.UUID]int)
	for _, paymentServiceItem := range oldPaymentRequest.PaymentServiceItems {
		count := oldMTOServiceItemIDs[paymentServiceItem.MTOServiceItemID]
		oldMTOServiceItemIDs[paymentServiceItem.MTOServiceItemID] = count + 1
	}
	newMTOServiceItemIDs := make(map[uuid.UUID]int)
	for _, paymentServiceItem := range newPaymentRequest.PaymentServiceItems {
		count := newMTOServiceItemIDs[paymentServiceItem.MTOServiceItemID]
		newMTOServiceItemIDs[paymentServiceItem.MTOServiceItemID] = count + 1
	}
	suite.Equal(oldMTOServiceItemIDs, newMTOServiceItemIDs, "Referenced MTOServiceItems are not the same")

	// Test the service items, prices, and expected changed parameters.  Note that we don't check
	// all parameters since we assume the payment request creator we're calling has already tested
	// that functionality.
	type paramMap struct {
		name  models.ServiceItemParamName
		value string
	}

	strTestOriginalWeight := strconv.Itoa(repriceTestOriginalWeight.Int())
	strTestChangedOriginalWeight := strconv.Itoa(repriceTestNewOriginalWeight.Int())
	testServicePriceParams := []struct {
		isNewPaymentRequest bool
		paymentRequest      *models.PaymentRequest
		serviceCode         models.ReServiceCode
		priceCents          unit.Cents
		paramsToCheck       []paramMap
	}{
		// Old payment request that we were repricing
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeMS,
			priceCents:     unit.Cents(25513),
		},
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeCS,
			priceCents:     unit.Cents(22399),
		},
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeDLH,
			priceCents:     unit.Cents(281402),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strTestOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strTestOriginalWeight},
			},
		},
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeFSC,
			priceCents:     unit.Cents(1420),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strTestOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strTestOriginalWeight},
			},
		},
		// New payment request with new prices
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeMS,
			priceCents:          unit.Cents(25513),
		},
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeCS,
			priceCents:          unit.Cents(22399),
		},
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeDLH,
			priceCents:          unit.Cents(262909),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strTestChangedOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strTestChangedOriginalWeight},
			},
		},
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeFSC,
			priceCents:          unit.Cents(1420), // Price same as before since new weight still in same weight bracket
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strTestChangedOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strTestChangedOriginalWeight},
			},
		},
	}

	for _, servicePriceParam := range testServicePriceParams {
		label := "for old payment request"
		if servicePriceParam.isNewPaymentRequest {
			label = "for new payment request"
		}
		foundService := false
		for _, paymentServiceItem := range servicePriceParam.paymentRequest.PaymentServiceItems {
			if paymentServiceItem.MTOServiceItem.ReService.Code == servicePriceParam.serviceCode {
				foundService = true
				if suite.NotNilf(paymentServiceItem.PriceCents, "Price should not be nil for service code %s (%s)", servicePriceParam.serviceCode, label) {
					suite.Equalf(servicePriceParam.priceCents, *paymentServiceItem.PriceCents, "Prices do not match for service code %s (%s)", servicePriceParam.serviceCode, label)
				}
				for _, paramToCheck := range servicePriceParam.paramsToCheck {
					foundParam := false
					for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
						if paymentServiceItemParam.ServiceItemParamKey.Key == paramToCheck.name {
							foundParam = true
							suite.Equal(paramToCheck.value, paymentServiceItemParam.Value)
							break
						}
					}
					suite.Truef(foundParam, "Could not find param %s for service code %s (%s)", paramToCheck.name, servicePriceParam.serviceCode, label)
				}
				break
			}
		}
		suite.Truef(foundService, "Could not find service code %s (%s)", servicePriceParam.serviceCode, label)
	}

	// Check the proof of service docs to make sure they have the same core data
	if suite.Equal(len(oldPaymentRequest.ProofOfServiceDocs), len(newPaymentRequest.ProofOfServiceDocs), "Both payment requests should have same number of proof of service docs") {
		for i := 0; i < len(oldPaymentRequest.ProofOfServiceDocs); i++ {
			suite.Equal(newPaymentRequest.ID, newPaymentRequest.ProofOfServiceDocs[i].PaymentRequestID, "Proof of service doc should point to the new payment request ID")
		}
	}
	oldUploadIDs := make(map[uuid.UUID]int)
	for _, proofOfServiceDoc := range oldPaymentRequest.ProofOfServiceDocs {
		for _, primeUpload := range proofOfServiceDoc.PrimeUploads {
			count := oldUploadIDs[primeUpload.UploadID]
			oldUploadIDs[primeUpload.UploadID] = count + 1
		}
	}
	newUploadIDs := make(map[uuid.UUID]int)
	for _, proofOfServiceDoc := range newPaymentRequest.ProofOfServiceDocs {
		for _, primeUpload := range proofOfServiceDoc.PrimeUploads {
			count := newUploadIDs[primeUpload.UploadID]
			newUploadIDs[primeUpload.UploadID] = count + 1
		}
	}
	suite.Equal(oldUploadIDs, newUploadIDs, "Referenced UploadIDs are not the same")

	// Make sure the links between payment requests are set up properly.
	suite.Nil(oldPaymentRequest.RepricedPaymentRequestID, "Old payment request should have nil link")
	if suite.NotNil(newPaymentRequest.RepricedPaymentRequestID, "New payment request should not have nil link") {
		suite.Equal(oldPaymentRequest.ID, *newPaymentRequest.RepricedPaymentRequestID, "New payment request should link to the old payment request ID")
	}
}

func (suite *PaymentRequestServiceSuite) TestRepricePaymentRequestErrors() {
	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("Zip3TransitDistance",
		repriceTestPickupZip,
		repriceTestDestinationZip,
	).Return(repriceTestZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	repricer := NewPaymentRequestRepricer(creator)

	suite.T().Run("Fail to find payment request ID", func(t *testing.T) {
		bogusPaymentRequestID := uuid.Must(uuid.NewV4())
		newPaymentRequest, err := repricer.RepricePaymentRequest(suite.TestAppContext(), bogusPaymentRequestID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.Contains(err.Error(), "no rows in result set")
		}
	})

	suite.T().Run("Old payment status has unexpected status", func(t *testing.T) {
		paidPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusPaid,
			},
		})
		newPaymentRequest, err := repricer.RepricePaymentRequest(suite.TestAppContext(), paidPaymentRequest.ID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.Contains(err.Error(), "only pending payment requests can be repriced, but this payment request has status of "+models.PaymentRequestStatusPaid)
		}
	})

	suite.T().Run("Can handle error when creating new repriced payment request", func(t *testing.T) {
		// Mock out a creator.
		mockCreator := &mocks.PaymentRequestCreator{}
		mockCreator.On("CreatePaymentRequest",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest"),
		).Return(nil, errors.New("test error"))

		repricerWithMockCreator := NewPaymentRequestRepricer(mockCreator)

		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		newPaymentRequest, err := repricerWithMockCreator.RepricePaymentRequest(suite.TestAppContext(), paymentRequest.ID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.Contains(err.Error(), "test error")
		}
	})
}

func (suite *PaymentRequestServiceSuite) setupRepriceData() (models.Move, models.PaymentRequest) {
	// Pickup/destination addresses
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "235 Prospect Valley Road SE",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     repriceTestPickupZip,
		},
	})
	destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "17 8th St",
			City:           "San Antonio",
			State:          "TX",
			PostalCode:     repriceTestDestinationZip,
		},
	})

	// Contract year, service area, rate area, zip3
	contractYear, serviceArea, _, _ := testdatagen.SetupServiceAreaRateArea(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EscalationCompounded: repriceTestEscalationCompounded,
		},
		ReRateArea: models.ReRateArea{
			Name: "Georgia",
		},
		ReZip3: models.ReZip3{
			Zip3:          pickupAddress.PostalCode[0:3],
			BasePointCity: pickupAddress.City,
			State:         pickupAddress.State,
		},
	})

	// MS price data
	msService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeMS,
		},
	})
	msTaskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      msService.ID,
		PriceCents:     repriceTestMSFee,
	}
	suite.MustSave(&msTaskOrderFee)

	// CS price data
	csService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeCS,
		},
	})
	csTaskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      csService.ID,
		PriceCents:     repriceTestCSFee,
	}
	suite.MustSave(&csTaskOrderFee)

	// DLH price data
	testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          false,
			PriceMillicents:       repriceTestDLHPrice,
		},
	})

	// Create move, shipment, and service items for MS, CS, DLH, and FSC.
	availableToPrimeAt := time.Date(testdatagen.GHCTestYear, time.July, 1, 0, 0, 0, 0, time.UTC)
	estimatedWeight := repriceTestEstimatedWeight
	originalWeight := repriceTestOriginalWeight
	moveTaskOrder, mtoServiceItems := testdatagen.MakeFullDLHMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: &availableToPrimeAt,
		},
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &originalWeight,
			PickupAddressID:      &pickupAddress.ID,
			PickupAddress:        &pickupAddress,
			DestinationAddressID: &destinationAddress.ID,
			DestinationAddress:   &destinationAddress,
		},
	})

	// FSC price data (needs actual pickup date from move created above)
	publicationDate := moveTaskOrder.MTOShipments[0].ActualPickupDate.AddDate(0, 0, -3) // 3 days earlier
	ghcDieselFuelPrice := models.GHCDieselFuelPrice{
		PublicationDate:       publicationDate,
		FuelPriceInMillicents: repriceTestFSCPrice,
	}
	suite.MustSave(&ghcDieselFuelPrice)

	// Build up a payment request with service item references for creating a payment request.
	paymentRequestArg := models.PaymentRequest{
		MoveTaskOrderID:     moveTaskOrder.ID,
		IsFinal:             false,
		PaymentServiceItems: models.PaymentServiceItems{},
	}
	for _, mtoServiceItem := range mtoServiceItems {
		newPaymentServiceItem := models.PaymentServiceItem{
			MTOServiceItemID: mtoServiceItem.ID,
			MTOServiceItem:   mtoServiceItem,
		}
		paymentRequestArg.PaymentServiceItems = append(paymentRequestArg.PaymentServiceItems, newPaymentServiceItem)
	}

	return moveTaskOrder, paymentRequestArg
}
