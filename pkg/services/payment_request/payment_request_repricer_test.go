package paymentrequest

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	repriceTestPickupZip             = "30907"
	repriceTestDestinationZip        = "78234"
	repriceTestMSFee                 = unit.Cents(25513)
	repriceTestCSFee                 = unit.Cents(22399)
	repriceTestDLHPrice              = unit.Millicents(6000)
	repriceTestFSCPrice              = unit.Millicents(277600)
	repriceTestEstimatedWeight       = unit.Pound(3500)
	repriceTestOriginalWeight        = unit.Pound(3652)
	repriceTestChangedOriginalWeight = unit.Pound(3412)
	repriceTestEscalationCompounded  = 1.04071
	repriceTestZip3Distance          = 1234
)

func (suite *PaymentRequestServiceSuite) TestRepricePaymentRequest() {
	// Setup baseline move/shipment/service items data along with needed rate data.
	move, paymentRequestArg := suite.setupRepriceData()

	// Mock out a planner.
	planner := &routemocks.Planner{}
	planner.On("Zip3TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(repriceTestZip3Distance, nil)

	// Create an initial payment request.
	creator := NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	paymentRequest, err := creator.CreatePaymentRequest(suite.TestAppContext(), &paymentRequestArg)
	suite.FatalNoError(err)

	// Adjust shipment's original weight to force different pricing on a reprice.
	mtoShipment := move.MTOShipments[0]
	newWeight := repriceTestChangedOriginalWeight
	mtoShipment.PrimeActualWeight = &newWeight
	suite.MustSave(&mtoShipment)

	// Reprice the payment request created above.
	repricer := NewPaymentRequestRepricer(creator)
	repricedPaymentRequest, err := repricer.RepricePaymentRequest(suite.TestAppContext(), paymentRequest.ID)
	suite.FatalNoError(err)

	// Fetch the old payment request again (since repricing should have changed its status).
	// Need to eager fetch some related data to use in test assertions below.
	var reloadedPaymentRequest models.PaymentRequest
	err = suite.DB().
		EagerPreload(
			"PaymentServiceItems.MTOServiceItem.ReService",
			"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey").
		Find(&reloadedPaymentRequest, paymentRequest.ID)
	suite.FatalNoError(err)

	// Verify some top-level items on the payment requests.
	suite.Equal(reloadedPaymentRequest.MoveTaskOrderID, repricedPaymentRequest.MoveTaskOrderID, "both payment requests should point to same move")
	suite.Equal(len(reloadedPaymentRequest.PaymentServiceItems), len(repricedPaymentRequest.PaymentServiceItems), "both payment requests should have same number of service items")
	suite.Equal(reloadedPaymentRequest.Status, models.PaymentRequestStatusReviewedAllRejected, "Exiting payment request status incorrect")
	suite.Equal(repricedPaymentRequest.Status, models.PaymentRequestStatusPending, "Repriced payment request status incorrect")

	// Verify that the IDs of the MTO service items remain the same across both payment requests.
	existingMTOServiceItems := make(map[uuid.UUID]int)
	for _, paymentServiceItem := range reloadedPaymentRequest.PaymentServiceItems {
		count := existingMTOServiceItems[paymentServiceItem.MTOServiceItemID]
		existingMTOServiceItems[paymentServiceItem.MTOServiceItemID] = count + 1
	}
	repricedMTOServiceItems := make(map[uuid.UUID]int)
	for _, paymentServiceItem := range repricedPaymentRequest.PaymentServiceItems {
		count := repricedMTOServiceItems[paymentServiceItem.MTOServiceItemID]
		repricedMTOServiceItems[paymentServiceItem.MTOServiceItemID] = count + 1
	}
	suite.Equal(existingMTOServiceItems, repricedMTOServiceItems, "Referenced MTOServiceItems are not the same")

	// Test the service items, prices, and expected changed parameters.  Note that we don't check
	// all parameters since we assume the payment request creator we're calling has already tested
	// that functionality.
	type paramMap struct {
		name  models.ServiceItemParamName
		value string
	}

	strRepriceTestOriginalWeight := strconv.Itoa(repriceTestOriginalWeight.Int())
	strRepriceTestChangedOriginalWeight := strconv.Itoa(repriceTestChangedOriginalWeight.Int())
	testServicePriceParams := []struct {
		isRepriced     bool
		paymentRequest *models.PaymentRequest
		serviceCode    models.ReServiceCode
		priceCents     unit.Cents
		paramsToCheck  []paramMap
	}{
		// Existing payment request that we were repricing
		{
			paymentRequest: &reloadedPaymentRequest,
			serviceCode:    models.ReServiceCodeMS,
			priceCents:     unit.Cents(25513),
		},
		{
			paymentRequest: &reloadedPaymentRequest,
			serviceCode:    models.ReServiceCodeCS,
			priceCents:     unit.Cents(22399),
		},
		{
			paymentRequest: &reloadedPaymentRequest,
			serviceCode:    models.ReServiceCodeDLH,
			priceCents:     unit.Cents(281402),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strRepriceTestOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strRepriceTestOriginalWeight},
			},
		},
		{
			paymentRequest: &reloadedPaymentRequest,
			serviceCode:    models.ReServiceCodeFSC,
			priceCents:     unit.Cents(1420),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strRepriceTestOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strRepriceTestOriginalWeight},
			},
		},
		// New payment request with new prices
		{
			isRepriced:     true,
			paymentRequest: repricedPaymentRequest,
			serviceCode:    models.ReServiceCodeMS,
			priceCents:     unit.Cents(25513),
		},
		{
			isRepriced:     true,
			paymentRequest: repricedPaymentRequest,
			serviceCode:    models.ReServiceCodeCS,
			priceCents:     unit.Cents(22399),
		},
		{
			isRepriced:     true,
			paymentRequest: repricedPaymentRequest,
			serviceCode:    models.ReServiceCodeDLH,
			priceCents:     unit.Cents(262909),
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strRepriceTestChangedOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strRepriceTestChangedOriginalWeight},
			},
		},
		{
			isRepriced:     true,
			paymentRequest: repricedPaymentRequest,
			serviceCode:    models.ReServiceCodeFSC,
			priceCents:     unit.Cents(1420), // Price same as before since new weight still in same weight bracket
			paramsToCheck: []paramMap{
				{models.ServiceItemParamNameWeightOriginal, strRepriceTestChangedOriginalWeight},
				{models.ServiceItemParamNameWeightBilled, strRepriceTestChangedOriginalWeight},
			},
		},
	}

	for _, servicePriceParam := range testServicePriceParams {
		label := "for existing payment request"
		if servicePriceParam.isRepriced {
			label = "for repriced payment request"
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
