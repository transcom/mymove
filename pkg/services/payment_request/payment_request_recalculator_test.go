package paymentrequest

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	recalculateTestPickupZip                 = "30907"
	recalculateTestDestinationZip            = "78234"
	recalculateTestMSFee                     = unit.Cents(12303)
	recalculateTestCSFee                     = unit.Cents(12303)
	recalculateTestDLHPrice                  = unit.Millicents(6000)
	recalculateTestFSCPrice                  = unit.Millicents(277600)
	recalculateTestDomOtherPrice             = unit.Cents(2159)
	recalculateTestDomServiceAreaPriceDOP    = unit.Cents(2359)
	recalculateTestDomServiceAreaPriceDOASIT = unit.Cents(335)
	recalculateTestEstimatedWeight           = unit.Pound(3500)
	recalculateTestOriginalWeight            = unit.Pound(3652)
	recalculateTestNewOriginalWeight         = unit.Pound(3412)
	recalculateTestEscalationCompounded      = 1.04071
	recalculateTestZip3Distance              = 1234
	recalculateNumProofOfServiceDocs         = 2
	recalculateNumberDaysSIT                 = 20
)

var (
	recalculateSITEntryDate           = time.Date(testdatagen.GHCTestYear, time.July, 15, 0, 0, 0, 0, time.UTC)
	recalculateSITPaymentRequestStart = recalculateSITEntryDate.AddDate(0, 0, 1).Format("2006-01-02")
	recalculateSITPaymentRequestEnd   = recalculateSITEntryDate.AddDate(0, 0, 20).Format("2006-01-02")
)

func (suite *PaymentRequestServiceSuite) TestRecalculatePaymentRequestSuccess() {
	// Setup baseline move/shipment/service items data along with needed rate data.
	move, paymentRequestArg := suite.setupRecalculateData1()

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	paymentRequest, err := creator.CreatePaymentRequestCheck(suite.AppContextForTest(), &paymentRequestArg)
	suite.FatalNoError(err)

	// Add a few proof of service docs and prime uploads.
	var oldProofOfServiceDocIDs []string
	for i := 0; i < recalculateNumProofOfServiceDocs; i++ {
		proofOfServiceDoc := factory.BuildProofOfServiceDoc(suite.DB(), []factory.Customization{
			{
				Model:    *paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		oldProofOfServiceDocIDs = append(oldProofOfServiceDocIDs, proofOfServiceDoc.ID.String())
		factory.BuildPrimeUpload(suite.DB(), []factory.Customization{
			{
				Model:    proofOfServiceDoc,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPrimeUpload(suite.DB(), []factory.Customization{
			{
				Model:    proofOfServiceDoc,
				LinkOnly: true,
			},
		}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})
	}
	sort.Strings(oldProofOfServiceDocIDs)

	// Adjust shipment's original weight to force different pricing on a recalculation.
	mtoShipment := move.MTOShipments[0]
	newWeight := recalculateTestNewOriginalWeight
	mtoShipment.PrimeActualWeight = &newWeight
	suite.MustSave(&mtoShipment)

	// Recalculate the payment request created above.
	statusUpdater := NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := NewPaymentRequestRecalculator(creator, statusUpdater)
	newPaymentRequest, err := recalculator.RecalculatePaymentRequest(suite.AppContextForTest(), paymentRequest.ID)
	suite.FatalNoError(err)

	// Fetch the old payment request again -- status should have changed and it should no longer
	// have proof of service docs now.  Need to eager fetch some related data to use in test
	// assertions below.
	var oldPaymentRequest models.PaymentRequest
	err = suite.DB().
		EagerPreload(
			"PaymentServiceItems.MTOServiceItem.ReService",
			"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
			"ProofOfServiceDocs",
		).
		Find(&oldPaymentRequest, paymentRequest.ID)
	suite.FatalNoError(err)

	// Verify some top-level items on the payment requests.
	suite.Equal(oldPaymentRequest.MoveTaskOrderID, newPaymentRequest.MoveTaskOrderID, "Both payment requests should point to same move")
	suite.Len(oldPaymentRequest.PaymentServiceItems, 5)
	suite.Equal(len(oldPaymentRequest.PaymentServiceItems), len(newPaymentRequest.PaymentServiceItems), "Both payment requests should have same number of service items")
	suite.Equal(oldPaymentRequest.Status, models.PaymentRequestStatusDeprecated, "Old payment request status incorrect")
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

	strTestOriginalWeight := strconv.Itoa(recalculateTestOriginalWeight.Int())
	strTestChangedOriginalWeight := strconv.Itoa(recalculateTestNewOriginalWeight.Int())
	strNumberDaysSIT := strconv.Itoa(recalculateNumberDaysSIT)
	testServicePriceParams := []struct {
		isNewPaymentRequest bool
		paymentRequest      *models.PaymentRequest
		serviceCode         models.ReServiceCode
		priceCents          unit.Cents
		paramsToCheck       []paramMap
	}{
		// Old payment request that we were recalculating
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeMS,
			priceCents:     unit.Cents(12303),
		},
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeCS,
			priceCents:     unit.Cents(12303),
		},
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeDLH,
			priceCents:     unit.Cents(279407),
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
		{
			paymentRequest: &oldPaymentRequest,
			serviceCode:    models.ReServiceCodeDOASIT,
			priceCents:     unit.Cents(254910),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strTestOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strTestOriginalWeight},
				{models.ServiceItemParamNameSITPaymentRequestStart, recalculateSITPaymentRequestStart},
				{models.ServiceItemParamNameSITPaymentRequestEnd, recalculateSITPaymentRequestEnd},
				{models.ServiceItemParamNameNumberDaysSIT, strNumberDaysSIT},
			},
		},
		// New payment request with new prices
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeMS,
			priceCents:          unit.Cents(12303),
		},
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeCS,
			priceCents:          unit.Cents(12303),
		},
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeDLH,
			priceCents:          unit.Cents(261045),
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
		{
			isNewPaymentRequest: true,
			paymentRequest:      newPaymentRequest,
			serviceCode:         models.ReServiceCodeDOASIT,
			priceCents:          unit.Cents(238158), // Price same as before since new weight still in same weight bracket
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strTestChangedOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strTestChangedOriginalWeight},
				{models.ServiceItemParamNameSITPaymentRequestStart, recalculateSITPaymentRequestStart},
				{models.ServiceItemParamNameSITPaymentRequestEnd, recalculateSITPaymentRequestEnd},
				{models.ServiceItemParamNameNumberDaysSIT, strNumberDaysSIT},
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

	// Check the proof of service docs; old payment request should have no proof of service docs now; new payment
	// request should have all the old payment request's proof of service docs.
	suite.Len(oldPaymentRequest.ProofOfServiceDocs, 0)
	var newProofOfServiceDocIDs []string
	if suite.Len(newPaymentRequest.ProofOfServiceDocs, recalculateNumProofOfServiceDocs) {
		for _, proofOfServiceDoc := range newPaymentRequest.ProofOfServiceDocs {
			suite.Equal(newPaymentRequest.ID, proofOfServiceDoc.PaymentRequestID, "Proof of service doc should point to the new payment request ID")
			newProofOfServiceDocIDs = append(newProofOfServiceDocIDs, proofOfServiceDoc.ID.String())
		}
	}
	sort.Strings(newProofOfServiceDocIDs)
	suite.Equal(oldProofOfServiceDocIDs, newProofOfServiceDocIDs, "Proof of service doc IDs differ, but should be the same")

	// Make sure the links between payment requests are set up properly.
	suite.Nil(oldPaymentRequest.RecalculationOfPaymentRequestID, "Old payment request should have nil link")
	if suite.NotNil(newPaymentRequest.RecalculationOfPaymentRequestID, "New payment request should not have nil link") {
		suite.Equal(oldPaymentRequest.ID, *newPaymentRequest.RecalculationOfPaymentRequestID, "New payment request should link to the old payment request ID")
	}
}

func (suite *PaymentRequestServiceSuite) TestRecalculatePaymentRequestErrors() {
	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := NewPaymentRequestRecalculator(creator, statusUpdater)

	suite.Run("Fail to find payment request ID", func() {
		bogusPaymentRequestID := uuid.Must(uuid.NewV4())
		newPaymentRequest, err := recalculator.RecalculatePaymentRequest(suite.AppContextForTest(), bogusPaymentRequestID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), bogusPaymentRequestID.String())
		}
	})

	suite.Run("Old payment status has unexpected status", func() {
		paidPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPaid,
				},
			},
		}, nil)
		newPaymentRequest, err := recalculator.RecalculatePaymentRequest(suite.AppContextForTest(), paidPaymentRequest.ID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.IsType(apperror.ConflictError{}, err)
			suite.Contains(err.Error(), paidPaymentRequest.ID.String())
			suite.Contains(err.Error(), models.PaymentRequestStatusPaid)
		}
	})

	suite.Run("Can handle error when creating new recalculated payment request", func() {
		// Mock out a creator.
		errString := "mock creator test error"
		mockCreator := &mocks.PaymentRequestCreator{}
		mockCreator.On("CreatePaymentRequestCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest"),
		).Return(nil, errors.New(errString))

		recalculatorWithMockCreator := NewPaymentRequestRecalculator(mockCreator, statusUpdater)

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		newPaymentRequest, err := recalculatorWithMockCreator.RecalculatePaymentRequest(suite.AppContextForTest(), paymentRequest.ID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.Equal(err.Error(), errString)
		}
	})

	suite.Run("Can handle error when updating old payment request status", func() {
		// Mock out a status updater.
		errString := "mock status updater test error"
		mockStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		mockStatusUpdater.On("UpdatePaymentRequestStatus",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PaymentRequest"),
			mock.AnythingOfType("string"),
		).Return(nil, errors.New(errString))

		recalculatorWithMockStatusUpdater := NewPaymentRequestRecalculator(creator, mockStatusUpdater)

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		newPaymentRequest, err := recalculatorWithMockStatusUpdater.RecalculatePaymentRequest(suite.AppContextForTest(), paymentRequest.ID)
		suite.Nil(newPaymentRequest)
		if suite.Error(err) {
			suite.Equal(err.Error(), errString)
		}
	})
}

func (suite *PaymentRequestServiceSuite) setupRecalculateData1() (models.Move, models.PaymentRequest) {
	// Pickup/destination addresses
	pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "235 Prospect Valley Road SE",
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     recalculateTestPickupZip,
			},
		},
	}, nil)
	destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "17 8th St",
				City:           "San Antonio",
				State:          "TX",
				PostalCode:     recalculateTestDestinationZip,
			},
		},
	}, nil)

	// Contract year, service area, rate area, zip3
	contractYear, serviceArea, _, _ := testdatagen.SetupServiceAreaRateArea(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EscalationCompounded: recalculateTestEscalationCompounded,
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

	// DLH price data
	testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          false,
			PriceMillicents:       recalculateTestDLHPrice,
		},
	})

	// Create move, shipment, and service items for MS, CS, DLH, and FSC.
	availableToPrimeAt := time.Date(testdatagen.GHCTestYear, time.July, 1, 0, 0, 0, 0, time.UTC)
	estimatedWeight := recalculateTestEstimatedWeight
	originalWeight := recalculateTestOriginalWeight
	moveTaskOrder, mtoServiceItems := factory.BuildFullDLHMTOServiceItems(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
				ApprovedAt:         &availableToPrimeAt,
			},
		},
		{
			Model:    pickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &originalWeight,
				SITDaysAllowance:     models.IntPointer(90),
			},
		},
	}, nil)

	// FSC price data (needs actual pickup date from move created above)
	publicationDate := moveTaskOrder.MTOShipments[0].ActualPickupDate.AddDate(0, 0, -3) // 3 days earlier
	effectiveDate := publicationDate.AddDate(0, 0, 1)
	endDate := publicationDate.AddDate(0, 0, 7)
	ghcDieselFuelPrice := models.GHCDieselFuelPrice{
		PublicationDate:       publicationDate,
		FuelPriceInMillicents: recalculateTestFSCPrice,
		EffectiveDate:         effectiveDate,
		EndDate:               endDate,
	}
	suite.MustSave(&ghcDieselFuelPrice)

	//  Domestic Origin Price Service
	domOriginPriceService := factory.BuildReServiceByCode(suite.DB(), models.ReServiceCodeDOP)

	domServiceAreaPriceDOP := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		ServiceID:             domOriginPriceService.ID,
		IsPeakPeriod:          false,
		Contract:              contractYear.Contract,
		DomesticServiceAreaID: serviceArea.ID,
		DomesticServiceArea:   serviceArea,
		PriceCents:            recalculateTestDomServiceAreaPriceDOP,
		Service:               domOriginPriceService,
	}
	suite.MustSave(&domServiceAreaPriceDOP)

	// Domestic Pack
	dpkService := factory.BuildReServiceByCode(suite.DB(), models.ReServiceCodeDPK)

	// Domestic Other Price
	domOtherPriceDPK := models.ReDomesticOtherPrice{
		ContractID:   contractYear.Contract.ID,
		ServiceID:    dpkService.ID,
		IsPeakPeriod: false,
		Schedule:     2,
		PriceCents:   recalculateTestDomOtherPrice,
		Contract:     contractYear.Contract,
		Service:      dpkService,
	}
	suite.MustSave(&domOtherPriceDPK)

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
			Status:           models.PaymentServiceItemStatusApproved,
		}
		paymentRequestArg.PaymentServiceItems = append(paymentRequestArg.PaymentServiceItems, newPaymentServiceItem)
	}

	// DOASIT
	mtoServiceItemDOASIT := factory.BuildRealMTOServiceItemWithAllDeps(suite.DB(), models.ReServiceCodeDOASIT, moveTaskOrder, moveTaskOrder.MTOShipments[0],
		[]factory.Customization{
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
			},
		}, nil)
	mtoServiceItemDOASIT.SITEntryDate = &recalculateSITEntryDate
	suite.MustSave(&mtoServiceItemDOASIT)

	domServiceAreaPriceDOASIT := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		ServiceID:             mtoServiceItemDOASIT.ReServiceID,
		IsPeakPeriod:          false,
		DomesticServiceAreaID: serviceArea.ID,
		PriceCents:            recalculateTestDomServiceAreaPriceDOASIT,
	}
	suite.MustSave(&domServiceAreaPriceDOASIT)

	doasitPaymentServiceItem := models.PaymentServiceItem{
		MTOServiceItemID: mtoServiceItemDOASIT.ID,
		MTOServiceItem:   mtoServiceItemDOASIT,
		Status:           models.PaymentServiceItemStatusApproved,
	}
	doasitPaymentServiceItem.PaymentServiceItemParams = models.PaymentServiceItemParams{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       recalculateSITPaymentRequestStart,
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       recalculateSITPaymentRequestEnd,
		},
	}
	paymentRequestArg.PaymentServiceItems = append(paymentRequestArg.PaymentServiceItems, doasitPaymentServiceItem)

	return moveTaskOrder, paymentRequestArg
}

func (suite *PaymentRequestServiceSuite) setupRecalculateData2(move models.Move, shipment models.MTOShipment) (models.Move, models.PaymentRequest) {

	moveTaskOrder, mtoServiceItems := factory.BuildFullOriginMTOServiceItems(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
	}, nil)

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
