package mtoshipment

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	servicesMocks "github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mt "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type approveShipmentSubtestData struct {
	appCtx                 appcontext.AppContext
	move                   models.Move
	planner                *mocks.Planner
	shipmentApprover       services.ShipmentApprover
	mockedShipmentApprover services.ShipmentApprover
	mockedShipmentRouter   *servicesMocks.ShipmentRouter
	reServiceCodes         []models.ReServiceCode
	moveWeights            services.MoveWeights
	mtoUpdater             services.MoveTaskOrderUpdater
}

// Creates data for the TestApproveShipment function
func (suite *MTOShipmentServiceSuite) createApproveShipmentSubtestData() (subtestData *approveShipmentSubtestData) {
	subtestData = &approveShipmentSubtestData{}

	subtestData.move = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
	suite.False(verrs.HasAny())
	suite.FatalNoError(err)

	// Let's also create a transit time object with a zero upper bound for weight (this can happen in the table).
	ghcDomesticTransitTime0LbsUpper := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     10001,
		WeightLbsUpper:     0,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	verrs, err = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime0LbsUpper)
	suite.False(verrs.HasAny())
	suite.FatalNoError(err)

	// Let's create service codes in the DB
	// INPK and IHPK don't need this since they're not truncated
	subtestData.reServiceCodes = []models.ReServiceCode{
		models.ReServiceCodeDLH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDPK,
		models.ReServiceCodeDUPK,
	}

	for _, serviceCode := range subtestData.reServiceCodes {
		factory.FetchReServiceByCode(suite.DB(), serviceCode)
	}

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &servicesMocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &servicesMocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	subtestData.mockedShipmentRouter = &servicesMocks.ShipmentRouter{}

	router := NewShipmentRouter()

	builder := query.NewQueryBuilder()
	waf := entitlements.NewWeightAllotmentFetcher()
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	ppmEstimator := &servicesMocks.PPMEstimator{}
	queryBuilder := query.NewQueryBuilder()

	siCreator := mtoserviceitem.NewMTOServiceItemCreator(
		planner,
		builder,
		moveRouter,
		ghcrateengine.NewDomesticUnpackPricer(),
		ghcrateengine.NewDomesticPackPricer(),
		ghcrateengine.NewDomesticLinehaulPricer(),
		ghcrateengine.NewDomesticShorthaulPricer(),
		ghcrateengine.NewDomesticOriginPricer(),
		ghcrateengine.NewDomesticDestinationPricer(),
		ghcrateengine.NewFuelSurchargePricer(),
		ghcrateengine.NewDomesticDestinationFirstDaySITPricer(),
		ghcrateengine.NewDomesticDestinationSITDeliveryPricer(),
		ghcrateengine.NewDomesticDestinationAdditionalDaysSITPricer(),
		ghcrateengine.NewDomesticDestinationSITFuelSurchargePricer(),
		ghcrateengine.NewDomesticOriginFirstDaySITPricer(),
		ghcrateengine.NewDomesticOriginSITPickupPricer(),
		ghcrateengine.NewDomesticOriginAdditionalDaysSITPricer(),
		ghcrateengine.NewDomesticOriginSITFuelSurchargePricer())
	subtestData.mtoUpdater = mt.NewMoveTaskOrderUpdater(
		queryBuilder,
		siCreator,
		moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), ppmEstimator,
	)
	subtestData.planner = &mocks.Planner{}
	mockSender := setUpMockNotificationSender()
	subtestData.moveWeights = moverouter.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)

	subtestData.shipmentApprover = NewShipmentApprover(router, siCreator, subtestData.planner, subtestData.moveWeights, subtestData.mtoUpdater, moveRouter)
	subtestData.mockedShipmentApprover = NewShipmentApprover(subtestData.mockedShipmentRouter, siCreator, subtestData.planner, subtestData.moveWeights, subtestData.mtoUpdater, moveRouter)
	subtestData.appCtx = suite.AppContextWithSessionForTest(&auth.Session{
		ApplicationName: auth.OfficeApp,
		OfficeUserID:    uuid.Must(uuid.NewV4()),
	})

	const (
		dopTestServiceArea = "123"
		dopTestWeight      = 1212
	)

	pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "7 Q St",
				City:           "Birmingham",
				State:          "KY",
				PostalCode:     "40356",
			},
		},
	}, nil)

	contractYear, serviceArea, _, _ := testdatagen.SetupServiceAreaRateArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea: dopTestServiceArea,
		},
		ReRateArea: models.ReRateArea{
			Name: "Alabama",
		},
		ReZip3: models.ReZip3{
			Zip3:          pickupAddress.PostalCode[0:3],
			BasePointCity: pickupAddress.City,
			State:         pickupAddress.State,
		},
	})

	baseLinehaulPrice := testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.ContractID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          false,
		},
	})

	_ = testdatagen.MakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			Contract:              contractYear.Contract,
			DomesticServiceAreaID: serviceArea.ID,
			DomesticServiceArea:   serviceArea,
			IsPeakPeriod:          true,
			PriceMillicents:       baseLinehaulPrice.PriceMillicents - 2500, // minus $0.025
		},
	})

	domesticOriginService := factory.FetchReService(suite.DB(), []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOP,
				Name: "Dom. Origin Price",
			},
		},
	}, nil)

	domesticOriginPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		ServiceID:             domesticOriginService.ID,
		IsPeakPeriod:          true,
		DomesticServiceAreaID: serviceArea.ID,
		PriceCents:            146,
	}

	domesticOriginPeakPrice := domesticOriginPrice
	domesticOriginPeakPrice.PriceCents = 146

	domesticOriginNonpeakPrice := domesticOriginPrice
	domesticOriginNonpeakPrice.IsPeakPeriod = false
	domesticOriginNonpeakPrice.PriceCents = 127

	return subtestData
}

func (suite *MTOShipmentServiceSuite) TestApproveShipment() {
	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &servicesMocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &servicesMocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	ppmEstimator := &servicesMocks.PPMEstimator{}
	queryBuilder := query.NewQueryBuilder()
	planner := &mocks.Planner{}
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(
		planner,
		queryBuilder,
		moveRouter,
		ghcrateengine.NewDomesticUnpackPricer(),
		ghcrateengine.NewDomesticPackPricer(),
		ghcrateengine.NewDomesticLinehaulPricer(),
		ghcrateengine.NewDomesticShorthaulPricer(),
		ghcrateengine.NewDomesticOriginPricer(),
		ghcrateengine.NewDomesticDestinationPricer(),
		ghcrateengine.NewFuelSurchargePricer(),
		ghcrateengine.NewDomesticDestinationFirstDaySITPricer(),
		ghcrateengine.NewDomesticDestinationSITDeliveryPricer(),
		ghcrateengine.NewDomesticDestinationAdditionalDaysSITPricer(),
		ghcrateengine.NewDomesticDestinationSITFuelSurchargePricer(),
		ghcrateengine.NewDomesticOriginFirstDaySITPricer(),
		ghcrateengine.NewDomesticOriginSITPickupPricer(),
		ghcrateengine.NewDomesticOriginAdditionalDaysSITPricer(),
		ghcrateengine.NewDomesticOriginSITFuelSurchargePricer())
	mtoUpdater := mt.NewMoveTaskOrderUpdater(
		queryBuilder,
		siCreator,
		moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), ppmEstimator,
	)
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	suite.Run("If the international mtoShipment is approved successfully it should create pre approved mtoServiceItems and should NOT update pricing without port data", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			}}, nil)

		// we need to get the usPostRegionCityIDs based off of the ZIP for the addresses
		pickupUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "50314")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Des Moines",
					State:              "IA",
					PostalCode:         "50314",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		suite.FatalNoError(err)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		pickupDate := time.Now().Add(24 * time.Hour)
		internationalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusSubmitted,
					MarketCode:           models.MarketCodeInternational,
					PrimeEstimatedWeight: models.PoundPointer(unit.Pound(4000)),
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					RequestedPickupDate:  &pickupDate,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)

		shipmentRouter := NewShipmentRouter()
		waf := entitlements.NewWeightAllotmentFetcher()
		mockSender := setUpMockNotificationSender()
		moveWeights := moverouter.NewMoveWeights(NewShipmentReweighRequester(mockSender), waf)
		var serviceItemCreator services.MTOServiceItemCreator
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(500, nil)

		// Approve international shipment
		shipmentApprover := NewShipmentApprover(shipmentRouter, serviceItemCreator, planner, moveWeights, mtoUpdater, moveRouter)
		_, err = shipmentApprover.ApproveShipment(appCtx, internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReserviceCodes := []models.ReServiceCode{
			models.ReServiceCodePOEFSC,
			models.ReServiceCodeISLH,
			models.ReServiceCodeIHPK,
			models.ReServiceCodeIHUPK,
		}

		suite.Equal(len(expectedReserviceCodes), len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			suite.True(slices.Contains(expectedReserviceCodes, actualReServiceCode))
			// we should have pricing data on all but the POEFSC since we don't have the port data yet
			if serviceItems[i].ReService.Code != models.ReServiceCodePOEFSC {
				suite.NotNil(serviceItems[i].PricingEstimate)
			} else if serviceItems[i].ReService.Code == models.ReServiceCodePOEFSC {
				suite.Nil(serviceItems[i].PricingEstimate)
			}
		}
	})

	suite.Run("Given international mtoShipment is approved successfully pre-approved mtoServiceItems are created NTS CONUS to OCONUS", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					FacilityName: *models.StringPointer("Test Storage Name"),
					Email:        models.StringPointer("old@email.com"),
					LotNumber:    models.StringPointer("Test lot number"),
					Phone:        models.StringPointer("555-555-5555"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99507",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		pickupDate := time.Now().Add(24 * time.Hour)
		internationalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:          models.MarketCodeInternational,
					Status:              models.MTOShipmentStatusSubmitted,
					ShipmentType:        models.MTOShipmentTypeHHGIntoNTS,
					RequestedPickupDate: &pickupDate,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)

		shipmentRouter := NewShipmentRouter()
		var serviceItemCreator services.MTOServiceItemCreator
		var planner route.Planner
		var moveWeights services.MoveWeights

		// Approve international shipment
		shipmentApprover := NewShipmentApprover(shipmentRouter, serviceItemCreator, planner, moveWeights, mtoUpdater, moveRouter)
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReserviceCodes := []models.ReServiceCode{
			models.ReServiceCodeISLH,
			models.ReServiceCodeINPK,
		}

		suite.Equal(len(expectedReserviceCodes), len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			suite.True(slices.Contains(expectedReserviceCodes, actualReServiceCode))
		}
	})

	suite.Run("Given international mtoShipment is approved successfully pre-approved mtoServiceItems are created NTS OCONUS to CONUS", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					FacilityName: *models.StringPointer("Test Storage Name"),
					Email:        models.StringPointer("old@email.com"),
					LotNumber:    models.StringPointer("Test lot number"),
					Phone:        models.StringPointer("555-555-5555"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		pickupDate := time.Now().Add(24 * time.Hour)
		internationalShipment := factory.BuildNTSShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99507",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:          models.MarketCodeInternational,
					Status:              models.MTOShipmentStatusSubmitted,
					ShipmentType:        models.MTOShipmentTypeHHGIntoNTS,
					RequestedPickupDate: &pickupDate,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)

		shipmentRouter := NewShipmentRouter()
		var serviceItemCreator services.MTOServiceItemCreator
		var planner route.Planner
		var moveWeights services.MoveWeights

		// Approve international shipment
		shipmentApprover := NewShipmentApprover(shipmentRouter, serviceItemCreator, planner, moveWeights, mtoUpdater, moveRouter)
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReserviceCodes := []models.ReServiceCode{
			models.ReServiceCodeISLH,
			models.ReServiceCodePODFSC,
			models.ReServiceCodeINPK,
		}

		suite.Equal(len(expectedReserviceCodes), len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			suite.True(slices.Contains(expectedReserviceCodes, actualReServiceCode))
		}
	})

	suite.Run("Given international mtoShipment is approved successfully pre-approved mtoServiceItems are created NTS-R CONUS to OCONUS", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					FacilityName: *models.StringPointer("Test Storage Name"),
					Email:        models.StringPointer("old@email.com"),
					LotNumber:    models.StringPointer("Test lot number"),
					Phone:        models.StringPointer("555-555-5555"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		internationalShipment := factory.BuildNTSRShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99507",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:   models.MarketCodeInternational,
					Status:       models.MTOShipmentStatusSubmitted,
					ShipmentType: models.MTOShipmentTypeHHGOutOfNTS,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)

		shipmentRouter := NewShipmentRouter()
		var serviceItemCreator services.MTOServiceItemCreator
		var planner route.Planner
		var moveWeights services.MoveWeights

		// Approve international shipment
		shipmentApprover := NewShipmentApprover(shipmentRouter, serviceItemCreator, planner, moveWeights, mtoUpdater, moveRouter)
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReserviceCodes := []models.ReServiceCode{
			models.ReServiceCodeISLH,
			models.ReServiceCodePOEFSC,
			models.ReServiceCodeIHUPK,
		}

		suite.Equal(len(expectedReserviceCodes), len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			suite.True(slices.Contains(expectedReserviceCodes, actualReServiceCode))
		}
	})

	suite.Run("Given international mtoShipment is approved successfully pre-approved mtoServiceItems are created NTS-R OCONUS to CONUS", func() {
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					FacilityName: *models.StringPointer("Test Storage Name"),
					Email:        models.StringPointer("old@email.com"),
					LotNumber:    models.StringPointer("Test lot number"),
					Phone:        models.StringPointer("555-555-5555"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99507",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		internationalShipment := factory.BuildNTSRShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:   models.MarketCodeInternational,
					Status:       models.MTOShipmentStatusSubmitted,
					ShipmentType: models.MTOShipmentTypeHHGOutOfNTS,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)

		shipmentRouter := NewShipmentRouter()
		var serviceItemCreator services.MTOServiceItemCreator
		var planner route.Planner
		var moveWeights services.MoveWeights

		// Approve international shipment
		shipmentApprover := NewShipmentApprover(shipmentRouter, serviceItemCreator, planner, moveWeights, mtoUpdater, moveRouter)
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReserviceCodes := []models.ReServiceCode{
			models.ReServiceCodeISLH,
			models.ReServiceCodePODFSC,
			models.ReServiceCodeIHUPK,
		}

		suite.Equal(len(expectedReserviceCodes), len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			suite.True(slices.Contains(expectedReserviceCodes, actualReServiceCode))
		}
	})

	suite.Run("If the mtoShipment is approved successfully it should create approved mtoServiceItems", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner
		estimatedWeight := unit.Pound(1212)
		tomorrow := time.Now().AddDate(0, 0, 1)

		shipmentForAutoApprove := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusSubmitted,
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &tomorrow,
				},
			},
		}, nil)
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(500, nil)

		preApprovalTime := time.Now()
		shipment, approverErr := approver.ApproveShipment(appCtx, shipmentForAutoApprove.ID, shipmentForAutoApproveEtag)

		suite.NoError(approverErr)
		suite.Equal(move.ID, shipment.MoveTaskOrderID)

		err := appCtx.DB().Find(&fetchedShipment, shipmentForAutoApprove.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusApproved, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)

		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipmentForAutoApprove.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		suite.Equal(6, len(serviceItems))

		// All ApprovedAt times for service items should be the same, so just get the first one
		// Test that service item was approved within a few seconds of the current time
		suite.Assertions.WithinDuration(preApprovalTime, *serviceItems[0].ApprovedAt, 2*time.Second)

		// If we've gotten the shipment updated and fetched it without error then we can inspect the
		// service items created as a side effect to see if they are approved.
		for i := range serviceItems {
			suite.Equal(models.MTOServiceItemStatusApproved, serviceItems[i].Status)
			suite.Equal(subtestData.reServiceCodes[i], serviceItems[i].ReService.Code)
		}
	})

	suite.Run("approves shipment of type PPM and loads PPMShipment association", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		shipmentForAutoApprove := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
		}, nil)
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.Shipment.UpdatedAt)

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		shipment, approverErr := approver.ApproveShipment(appCtx, shipmentForAutoApprove.Shipment.ID, shipmentForAutoApproveEtag)

		suite.NoError(approverErr)
		suite.Equal(move.ID, shipment.MoveTaskOrderID)

		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
		suite.Equal(shipment.ID, shipmentForAutoApprove.Shipment.ID)

		suite.Equal(shipmentForAutoApprove.ID, shipment.PPMShipment.ID)
		suite.Equal(models.PPMShipmentStatusSubmitted, shipment.PPMShipment.Status)
	})

	suite.Run("If we act on a shipment with a weight that has a 0 upper weight it should still work", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(11000)
		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		shipmentHeavy := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
					RequestedPickupDate:  &tomorrow,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipmentHeavy.ID)
		suite.FatalNoError(err)
		err = suite.DB().Load(&createdShipment)
		suite.FatalNoError(err)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			createdShipment.PickupAddress.PostalCode,
			createdShipment.DestinationAddress.PostalCode,
		).Return(500, nil)

		shipmentHeavyEtag := etag.GenerateEtag(shipmentHeavy.UpdatedAt)
		_, err = approver.ApproveShipment(appCtx, shipmentHeavy.ID, shipmentHeavyEtag)
		suite.NoError(err)

		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipmentHeavy.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.Run("When status transition is not allowed, returns a ConflictStatusError", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover

		rejectionReason := "a reason"
		pickupDate := time.Now().Add(24 * time.Hour)
		rejectedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusRejected,
					RejectionReason:     &rejectionReason,
					RequestedPickupDate: &pickupDate,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		_, err := approver.ApproveShipment(appCtx, rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover

		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		_, err := approver.ApproveShipment(appCtx, staleShipment.ID, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		approver := subtestData.shipmentApprover

		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := approver.ApproveShipment(appCtx, badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("It calls Approve on the ShipmentRouter", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.mockedShipmentApprover
		shipmentRouter := subtestData.mockedShipmentRouter

		pickupDate := time.Now().Add(24 * time.Hour)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: &pickupDate,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)
		err = suite.DB().Load(&createdShipment, "MoveTaskOrder", "PickupAddress", "DestinationAddress")
		suite.FatalNoError(err)

		shipmentRouter.On("Approve", mock.AnythingOfType("*appcontext.appContext"), &createdShipment).Return(nil)

		_, err = approver.ApproveShipment(appCtx, shipment.ID, eTag)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(suite.T(), "Approve", 1)
	})

	suite.Run("If the mtoShipment uses external vendor not allowed to approve shipment", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		shipmentForAutoApprove := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:             models.MTOShipmentStatusSubmitted,
					UsesExternalVendor: true,
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
				},
			},
		}, nil)
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		shipment, approverErr := approver.ApproveShipment(appCtx, shipmentForAutoApprove.ID, shipmentForAutoApproveEtag)

		suite.Contains(approverErr.Error(), "shipment uses external vendor, cannot be approved for GHC Prime")
		suite.Equal(uuid.UUID{}, shipment.ID)

		err := appCtx.DB().Find(&fetchedShipment, shipmentForAutoApprove.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusSubmitted, fetchedShipment.Status)
		suite.Nil(shipment.ApprovedDate)
		suite.Nil(fetchedShipment.ApprovedDate)

		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipmentForAutoApprove.ID).All(&serviceItems)
		suite.NoError(err)

		suite.Equal(0, len(serviceItems))
	})

	suite.Run("Test that correct addresses are being used to calculate required delivery date", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
			models.ReServiceCodeDNPK,
		}

		for _, serviceCode := range expectedReServiceCodes {
			factory.FetchReServiceByCode(appCtx.DB(), serviceCode)
		}

		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// MTOShipmentTypeHHG: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// MTOShipmentTypeHHGIntoNTS: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, StorageFacility
		// MTOShipmentTypeHHGOutOfNTS: ScheduledPickupDate, NTSRecordedWeight, StorageFacility, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(1400)

		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

		pickupDate := time.Now().Add(24 * time.Hour)
		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
					RequestedPickupDate:  &pickupDate,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		ntsShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTS,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
					RequestedPickupDate:  &pickupDate,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		ntsrShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeHHGOutOfNTS,
					ScheduledPickupDate: &testdatagen.DateInsidePeakRateCycle,
					NTSRecordedWeight:   &estimatedWeight,
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: &pickupDate,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		var TransitDistancePickupArg string
		var TransitDistanceDestinationArg string

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil).Run(func(args mock.Arguments) {
			TransitDistancePickupArg = args.Get(1).(string)
			TransitDistanceDestinationArg = args.Get(2).(string)
		})

		testCases := []struct {
			shipment            models.MTOShipment
			pickupLocation      *models.Address
			destinationLocation *models.Address
		}{
			{hhgShipment, hhgShipment.PickupAddress, hhgShipment.DestinationAddress},
			{ntsShipment, ntsShipment.PickupAddress, &ntsShipment.StorageFacility.Address},
			{ntsrShipment, &ntsrShipment.StorageFacility.Address, ntsrShipment.DestinationAddress},
		}

		for _, testCase := range testCases {
			shipmentEtag := etag.GenerateEtag(testCase.shipment.UpdatedAt)
			_, err := approver.ApproveShipment(appCtx, testCase.shipment.ID, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, testCase.shipment.ID)
			suite.NoError(err)
			// We also should have a required delivery date
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			// Check that TransitDistance was called with the correct addresses
			suite.Equal(testCase.pickupLocation.PostalCode, TransitDistancePickupArg)
			suite.Equal(testCase.destinationLocation.PostalCode, TransitDistanceDestinationArg)
		}
	})
	suite.Run("Approval of a shipment with an estimated weight will update authorized weight", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		planner := subtestData.planner
		approver := subtestData.shipmentApprover
		estimatedWeight := unit.Pound(1234)
		pickupDate := time.Now().Add(24 * time.Hour)
		shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusSubmitted,
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &pickupDate,
				},
			},
		}, nil)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil)

		suite.Equal(8000, *shipment.MoveTaskOrder.Orders.Entitlement.AuthorizedWeight())

		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)

		_, approverErr := approver.ApproveShipment(appCtx, shipment.ID, shipmentEtag)
		suite.NoError(approverErr)

		err := appCtx.DB().Reload(shipment.MoveTaskOrder.Orders.Entitlement)
		suite.NoError(err)

		estimatedWeight110 := int(math.Round(float64(*shipment.PrimeEstimatedWeight) * 1.10))
		suite.Equal(estimatedWeight110, *shipment.MoveTaskOrder.Orders.Entitlement.AuthorizedWeight())
	})

	suite.Run("Approval of a shipment that exceeds excess weight will flag for excess weight", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		planner := subtestData.planner
		approver := subtestData.shipmentApprover
		estimatedWeight := unit.Pound(100000)
		pickupDate := time.Now().Add(24 * time.Hour)
		shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusSubmitted,
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &pickupDate,
				},
			},
		}, nil)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil)

		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)

		_, approverErr := approver.ApproveShipment(appCtx, shipment.ID, shipmentEtag)
		suite.NoError(approverErr)

		err := appCtx.DB().Reload(&shipment.MoveTaskOrder)
		suite.NoError(err)

		suite.NotNil(shipment.MoveTaskOrder.ExcessWeightQualifiedAt)
	})

	suite.Run("If the CONUS to OCONUS UB mtoShipment is approved successfully it should create pre approved mtoServiceItems", func() {
		internationalShipment := factory.BuildMTOShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:          models.MarketCodeInternational,
					Status:              models.MTOShipmentStatusSubmitted,
					ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
					RequestedPickupDate: &tomorrow,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)
		shipmentApprover := suite.createApproveShipmentSubtestData().shipmentApprover
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeUBP,
			models.ReServiceCodePOEFSC,
			models.ReServiceCodeIUBPK,
			models.ReServiceCodeIUBUPK,
		}
		expectedReServiceNames := []string{
			"International UB price",
			"International POE fuel surcharge",
			"International UB pack",
			"International UB unpack",
		}

		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			actualReServiceName := serviceItems[i].ReService.Name
			suite.True(slices.Contains(expectedReServiceCodes, actualReServiceCode), "Contains unexpected code: "+actualReServiceCode.String())
			suite.True(slices.Contains(expectedReServiceNames, actualReServiceName), "Contains unexpected name: "+actualReServiceName)
		}
	})

	suite.Run("If the OCONUS to CONUS UB mtoShipment is approved successfully it should create pre approved mtoServiceItems", func() {
		var scheduledPickupDate time.Time
		internationalShipment := factory.BuildMTOShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:          models.MarketCodeInternational,
					Status:              models.MTOShipmentStatusSubmitted,
					ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
					ScheduledPickupDate: &scheduledPickupDate,
					RequestedPickupDate: &tomorrow,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)
		shipmentApprover := suite.createApproveShipmentSubtestData().shipmentApprover
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeUBP,
			models.ReServiceCodePODFSC,
			models.ReServiceCodeIUBPK,
			models.ReServiceCodeIUBUPK,
		}
		expectedReServiceNames := []string{
			"International UB price",
			"International POD fuel surcharge",
			"International UB pack",
			"International UB unpack",
		}

		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			actualReServiceName := serviceItems[i].ReService.Name
			suite.True(slices.Contains(expectedReServiceCodes, actualReServiceCode), "Contains unexpected code: "+actualReServiceCode.String())
			suite.True(slices.Contains(expectedReServiceNames, actualReServiceName), "Contains unexpected name: "+actualReServiceName)
		}
	})

	suite.Run("If the OCONUS to OCONUS UB mtoShipment is approved successfully it should create pre approved mtoServiceItems", func() {
		var scheduledPickupDate time.Time
		internationalShipment := factory.BuildMTOShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode:          models.MarketCodeInternational,
					Status:              models.MTOShipmentStatusSubmitted,
					ShipmentType:        models.MTOShipmentTypeUnaccompaniedBaggage,
					ScheduledPickupDate: &scheduledPickupDate,
					RequestedPickupDate: &tomorrow,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Fairbanks",
					State:          "AK",
					PostalCode:     "99701",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		internationalShipmentEtag := etag.GenerateEtag(internationalShipment.UpdatedAt)
		shipmentApprover := suite.createApproveShipmentSubtestData().shipmentApprover
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), internationalShipment.ID, internationalShipmentEtag)
		suite.NoError(err)

		// Get created pre approved service items
		var serviceItems []models.MTOServiceItem
		err2 := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", internationalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err2)

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeUBP,
			models.ReServiceCodeIUBPK,
			models.ReServiceCodeIUBUPK,
		}
		expectedReServiceNames := []string{
			"International UB price",
			"International UB pack",
			"International UB unpack",
		}

		suite.Equal(3, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			actualReServiceCode := serviceItems[i].ReService.Code
			actualReServiceName := serviceItems[i].ReService.Name
			suite.True(slices.Contains(expectedReServiceCodes, actualReServiceCode), "Contains unexpected code: "+actualReServiceCode.String())
			suite.True(slices.Contains(expectedReServiceNames, actualReServiceName), "Contains unexpected name: "+actualReServiceName)
		}
	})

	suite.Run("Given invalid shipment error returned", func() {
		invalidShipment := factory.BuildMTOShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		invalidShipmentEtag := etag.GenerateEtag(invalidShipment.UpdatedAt)

		shipmentRouter := NewShipmentRouter()
		var serviceItemCreator services.MTOServiceItemCreator
		var planner route.Planner
		var moveWeights services.MoveWeights

		// Approve international shipment
		shipmentApprover := NewShipmentApprover(shipmentRouter, serviceItemCreator, planner, moveWeights, mtoUpdater, moveRouter)
		_, err := shipmentApprover.ApproveShipment(suite.AppContextForTest(), invalidShipment.ID, invalidShipmentEtag)
		suite.Error(err)
	})
}

func (suite *MTOShipmentServiceSuite) TestApproveShipmentValidation() {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	suite.Run("RequestedPickupDate validation check - must be in the future for shipment types other than PPM", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover

		testCases := []struct {
			input        *time.Time
			shipmentType models.MTOShipmentType
			shouldError  bool
		}{
			// HHG - Domestic
			{nil, models.MTOShipmentTypeHHG, true},
			{&time.Time{}, models.MTOShipmentTypeHHG, true},
			{&yesterday, models.MTOShipmentTypeHHG, true},
			{&now, models.MTOShipmentTypeHHG, true},
			{&tomorrow, models.MTOShipmentTypeHHG, false},
			// NTSR - RequestedPickupDate NOT Required
			{nil, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&time.Time{}, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&yesterday, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&now, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&tomorrow, models.MTOShipmentTypeHHGOutOfNTS, false},
			// UB - International
			{nil, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&time.Time{}, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&yesterday, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&now, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&tomorrow, models.MTOShipmentTypeUnaccompaniedBaggage, false},
			// PPM - should always pass validation
			{nil, models.MTOShipmentTypePPM, false},
			{&time.Time{}, models.MTOShipmentTypePPM, false},
			{&yesterday, models.MTOShipmentTypePPM, false},
			{&now, models.MTOShipmentTypePPM, false},
			{&tomorrow, models.MTOShipmentTypePPM, false},
		}

		for _, testCase := range testCases {
			// Default is HHG, but we set it explicitly below via the test cases
			var shipment models.MTOShipment
			if testCase.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				ghcDomesticTransitTime := models.GHCDomesticTransitTime{
					MaxDaysTransitTime: 12,
					WeightLbsLower:     0,
					WeightLbsUpper:     10000,
					DistanceMilesLower: 0,
					DistanceMilesUpper: 10000,
				}
				verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
				suite.Assert().False(verrs.HasAny())
				suite.NoError(err)
				moveForPrime := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
				shipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    moveForPrime,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:         testCase.shipmentType,
							RequestedPickupDate:  testCase.input,
							ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
							PrimeEstimatedWeight: models.PoundPointer(unit.Pound(4000)),
							Status:               models.MTOShipmentStatusSubmitted,
						},
					},
				}, nil)
			} else {
				shipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
					{
						Model:    move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        testCase.shipmentType,
							Status:              models.MTOShipmentStatusSubmitted,
							RequestedPickupDate: testCase.input,
						},
					},
				}, nil)
			}

			if testCase.input == nil || testCase.input.IsZero() {
				// Overwrite factory merge issue with customizations
				err := suite.DB().Q().RawQuery("UPDATE mto_shipments SET requested_pickup_date=? WHERE id=?", testCase.input, shipment.ID).Exec()
				suite.NoError(err)
			}

			eTag := etag.GenerateEtag(shipment.UpdatedAt)

			createdShipment := models.MTOShipment{}
			err := suite.DB().Find(&createdShipment, shipment.ID)
			suite.FatalNoError(err)
			err = suite.DB().Load(&createdShipment, "MoveTaskOrder", "PickupAddress", "DestinationAddress")
			suite.FatalNoError(err)

			_, err = approver.ApproveShipment(appCtx, shipment.ID, eTag)

			testCaseInputString := ""
			if testCase.input == nil {
				testCaseInputString = "nil"
			} else {
				testCaseInputString = (*testCase.input).String()
			}

			if testCase.shouldError {
				suite.Error(err)
				if testCase.input != nil && !(*testCase.input).IsZero() {
					suite.Equal("RequestedPickupDate must be greater than or equal to tomorrow's date.", err.Error())
				} else {
					suite.Contains(err.Error(), fmt.Sprintf("RequestedPickupDate is required to create or modify %s %s shipment", GetAorAnByShipmentType(testCase.shipmentType), testCase.shipmentType))
				}
			} else {
				suite.NoError(err, "Should not error for %s | %s", testCase.shipmentType, testCaseInputString)
			}
		}
	})
}

func (suite *MTOShipmentServiceSuite) TestApproveShipments() {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	suite.Run("Successfully approves multiple shipments", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		shipmentApprover := subtestData.shipmentApprover

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		shipment1 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: &tomorrow,
				},
			},
		}, nil)

		shipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: &tomorrow,
				},
			},
		}, nil)

		eTag1 := etag.GenerateEtag(shipment1.UpdatedAt)
		eTag2 := etag.GenerateEtag(shipment2.UpdatedAt)

		shipmentIdWithEtagArr := []services.ShipmentIdWithEtag{
			{
				ShipmentID: shipment1.ID,
				ETag:       eTag1,
			},
			{
				ShipmentID: shipment2.ID,
				ETag:       eTag2,
			},
		}
		approvedShipments, err := shipmentApprover.ApproveShipments(suite.AppContextForTest(), shipmentIdWithEtagArr)

		suite.NoError(err)
		suite.NotNil(approvedShipments)
		suite.Len(*approvedShipments, 2)
		suite.Equal(shipment1.ID, (*approvedShipments)[0].ID)
		suite.Equal(shipment2.ID, (*approvedShipments)[1].ID)
	})

	suite.Run("Returns error if one shipment approval fails", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		shipmentApprover := subtestData.shipmentApprover

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		shipment1 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: &tomorrow,
				},
			},
		}, nil)

		shipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:              models.MTOShipmentStatusSubmitted,
					RequestedPickupDate: &tomorrow,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment2.UpdatedAt)

		shipmentIdWithEtagArr := []services.ShipmentIdWithEtag{
			{
				ShipmentID: shipment1.ID,
				ETag:       eTag,
			},
			{
				ShipmentID: shipment2.ID,
				ETag:       eTag,
			},
		}

		approvedShipments, err := shipmentApprover.ApproveShipments(suite.AppContextForTest(), shipmentIdWithEtagArr)

		suite.Error(err)
		suite.Len(*approvedShipments, 0)
	})

	suite.Run("Given invalid shipment error returned", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		shipmentApprover := subtestData.shipmentApprover
		invalidShipment := factory.BuildMTOShipment(suite.AppContextForTest().DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
		}, nil)
		invalidShipmentEtag := etag.GenerateEtag(invalidShipment.UpdatedAt)

		shipmentIdWithEtagArr := []services.ShipmentIdWithEtag{
			{
				ShipmentID: invalidShipment.ID,
				ETag:       invalidShipmentEtag,
			},
		}

		_, err := shipmentApprover.ApproveShipments(suite.AppContextForTest(), shipmentIdWithEtagArr)
		suite.Error(err)
	})

	suite.Run("RequestedPickupDate validation check - must be in the future for shipment types other than PPM", func() {

		subtestData := suite.createApproveShipmentSubtestData()
		shipmentApprover := subtestData.shipmentApprover
		appCtx := subtestData.appCtx
		move := subtestData.move

		testCases := []struct {
			input        *time.Time
			shipmentType models.MTOShipmentType
			shouldError  bool
		}{
			// HHG - Domestic
			{nil, models.MTOShipmentTypeHHG, true},
			{&time.Time{}, models.MTOShipmentTypeHHG, true},
			{&yesterday, models.MTOShipmentTypeHHG, true},
			{&now, models.MTOShipmentTypeHHG, true},
			{&tomorrow, models.MTOShipmentTypeHHG, false},
			// NTSR - RequestedPickupDate NOT Required
			{nil, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&time.Time{}, models.MTOShipmentTypeHHGOutOfNTS, false},
			{&yesterday, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&now, models.MTOShipmentTypeHHGOutOfNTS, true},
			{&tomorrow, models.MTOShipmentTypeHHGOutOfNTS, false},
			// UB - International
			{nil, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&time.Time{}, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&yesterday, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&now, models.MTOShipmentTypeUnaccompaniedBaggage, true},
			{&tomorrow, models.MTOShipmentTypeUnaccompaniedBaggage, false},
			// PPM - should always pass validation
			{nil, models.MTOShipmentTypePPM, false},
			{&time.Time{}, models.MTOShipmentTypePPM, false},
			{&yesterday, models.MTOShipmentTypePPM, false},
			{&now, models.MTOShipmentTypePPM, false},
			{&tomorrow, models.MTOShipmentTypePPM, false},
		}

		for _, testCase := range testCases {
			// Default is HHG, but we set it explicitly below via the test cases
			var shipment models.MTOShipment
			var shipment2 models.MTOShipment
			if testCase.shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
				ghcDomesticTransitTime := models.GHCDomesticTransitTime{
					MaxDaysTransitTime: 12,
					WeightLbsLower:     0,
					WeightLbsUpper:     10000,
					DistanceMilesLower: 0,
					DistanceMilesUpper: 10000,
				}
				verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
				suite.Assert().False(verrs.HasAny())
				suite.NoError(err)
				moveForPrime := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
				shipment = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    moveForPrime,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:         testCase.shipmentType,
							RequestedPickupDate:  testCase.input,
							ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
							PrimeEstimatedWeight: models.PoundPointer(unit.Pound(4000)),
							Status:               models.MTOShipmentStatusSubmitted,
						},
					},
				}, nil)
				shipment2 = factory.BuildUBShipment(suite.DB(), []factory.Customization{
					{
						Model:    moveForPrime,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:         testCase.shipmentType,
							RequestedPickupDate:  testCase.input,
							ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
							PrimeEstimatedWeight: models.PoundPointer(unit.Pound(4000)),
							Status:               models.MTOShipmentStatusSubmitted,
						},
					},
				}, nil)
			} else {
				shipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
					{
						Model:    move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        testCase.shipmentType,
							Status:              models.MTOShipmentStatusSubmitted,
							RequestedPickupDate: testCase.input,
						},
					},
				}, nil)
				shipment2 = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
					{
						Model:    move,
						LinkOnly: true,
					},
					{
						Model: models.MTOShipment{
							ShipmentType:        testCase.shipmentType,
							Status:              models.MTOShipmentStatusSubmitted,
							RequestedPickupDate: testCase.input,
						},
					},
				}, nil)
			}

			if testCase.input == nil || testCase.input.IsZero() {
				// Overwrite factory merge issue with customizations
				err := suite.DB().Q().RawQuery("UPDATE mto_shipments SET requested_pickup_date=? WHERE id in (?,?)", testCase.input, shipment.ID, shipment2.ID).Exec()
				suite.NoError(err)
			}

			eTag1 := etag.GenerateEtag(shipment.UpdatedAt)
			eTag2 := etag.GenerateEtag(shipment2.UpdatedAt)

			shipmentIdWithEtagArr := []services.ShipmentIdWithEtag{
				{
					ShipmentID: shipment.ID,
					ETag:       eTag1,
				},
				{
					ShipmentID: shipment2.ID,
					ETag:       eTag2,
				},
			}

			approvedShipments, err := shipmentApprover.ApproveShipments(appCtx, shipmentIdWithEtagArr)

			testCaseInputString := ""
			if testCase.input == nil {
				testCaseInputString = "nil"
			} else {
				testCaseInputString = (*testCase.input).String()
			}

			if testCase.shouldError {
				suite.NotNil(approvedShipments, "Should return even with error for %s | %s", testCase.shipmentType, testCaseInputString)
				suite.Len(*approvedShipments, 0)
				suite.Error(err)
				if testCase.input != nil && !(*testCase.input).IsZero() {
					suite.Equal("RequestedPickupDate must be greater than or equal to tomorrow's date.", err.Error())
				} else {
					suite.Contains(err.Error(), fmt.Sprintf("RequestedPickupDate is required to create or modify %s %s shipment", GetAorAnByShipmentType(testCase.shipmentType), testCase.shipmentType))
				}
			} else {
				suite.NoError(err, "Should not error for %s | %s", testCase.shipmentType, testCaseInputString)
				suite.Len(*approvedShipments, 2)
				suite.NotNil(shipment)
			}
		}
	})
}

// Calculate final estimated price for NTS INPK
// the base price cents comes from IHPK because iHHG -> iNTS
func computeINPKExpectedPriceCents(
	basePriceCents int,
	escalationFactor float64,
	marketFactor float64,
	primeEstimatedWeightLbs int,
) unit.Cents {
	// Convert incoming cents to dollars just like how the calculate_escalated_price proc does before multiplying it
	// As well as add double rounding to the dollar amounts since the proc uses dollars instead of cents
	unroundedEscDollarAmount := float64(basePriceCents) / 100.0
	roundedEscDollars := math.Round(unroundedEscDollarAmount*escalationFactor*100) / 100 // rounds to 2 decimals
	cwt := (float64(primeEstimatedWeightLbs) * 1.1) / 100.0
	finalDollars := roundedEscDollars * marketFactor * cwt
	// Second rounding
	final := math.Round(finalDollars * 100)

	return unit.Cents(final)
}

func (suite *MTOShipmentServiceSuite) TestApproveShipmentBasicServiceItemEstimatePricing() {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	setupOconusToConusNtsShipment := func(estimatedWeight *unit.Pound) (models.StorageFacility, models.Address, models.Address, models.MTOShipment) {
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.StorageFacility{
					FacilityName: *models.StringPointer("Test Storage Name"),
					Email:        models.StringPointer("old@email.com"),
					LotNumber:    models.StringPointer("Test lot number"),
					Phone:        models.StringPointer("555-555-5555"),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99507",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "148 S East St",
					City:           "Miami",
					State:          "FL",
					PostalCode:     "94535",
				},
			},
		}, nil)

		shipment := factory.BuildNTSShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:             models.MoveStatusAPPROVED,
					AvailableToPrimeAt: &now,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MarketCode:           models.MarketCodeInternational,
					Status:               models.MTOShipmentStatusSubmitted,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTS,
					PrimeEstimatedWeight: estimatedWeight,
					RequestedPickupDate:  &tomorrow,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		return storageFacility, pickupAddress, destinationAddress, shipment
	}

	suite.Run("NTS OCONUS to CONUS estimate prices on approval", func() {
		subtestData := suite.createApproveShipmentSubtestData()

		testCases := []struct {
			name            string
			estimatedWeight *unit.Pound
		}{
			{
				name:            "Successfully applies estimated weight to INPK when estimated weight is present",
				estimatedWeight: models.PoundPointer(100000),
			},
			{
				name:            "Leaves estimated INPK price as nil when no estimated weight is present",
				estimatedWeight: nil,
			},
		}
		for _, tc := range testCases {
			appCtx := subtestData.appCtx
			planner := subtestData.planner
			approver := subtestData.shipmentApprover
			_, _, _, shipment := setupOconusToConusNtsShipment(tc.estimatedWeight)
			contract, err := models.FetchContractForMove(suite.AppContextForTest(), shipment.MoveTaskOrderID)
			suite.FatalNoError(err)

			planner.On("ZipTransitDistance",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(500, nil)

			// Aprove the shipment to trigger the estimate pricing proc on INPK
			shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
			_, approverErr := approver.ApproveShipment(appCtx, shipment.ID, shipmentEtag)
			suite.FatalNoError(approverErr)

			// Get created pre approved service items
			var serviceItems []models.MTOServiceItem
			err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
			suite.NoError(err)

			// Get the contract escalation factor
			var escalationFactor float64
			err = suite.DB().RawQuery(`
			SELECT calculate_escalation_factor($1, $2)
		`, contract.ID, shipment.RequestedPickupDate).First(&escalationFactor)
			suite.FatalNoError(err)

			// Verify our non-truncated escalation factor db value is as expected
			// this also tests the calculate_escalation_factor proc
			// This information was pulled from the migration scripts (Or just run db fresh and perform the lookups
			// manually, whichever is your cup of tea)
			suite.Equal(escalationFactor, 1.11)

			// Fetch the INPK market factor from the DB
			inpkReService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeINPK)
			ntsMarketFactor, err := models.FetchMarketFactor(suite.AppContextForTest(), contract.ID, inpkReService.ID, "O")
			suite.FatalNoError(err)

			// Assert basic service items
			// Basic = created immediately on their own, not requested
			// Accessorial = created at a later date if requested
			expectedServiceItems := map[models.ReServiceCode]*unit.Cents{
				// Not testing ISLH or PODFSC so their prices will be nil
				models.ReServiceCodeISLH:   nil,
				models.ReServiceCodePODFSC: nil,
				// Remember that we pass in IHPK base price, not INPK base price. INPK doesn't have a base price
				// because it uses IHPK for iHHG -> iNTS packing
				models.ReServiceCodeINPK: func() *unit.Cents {
					// Handle test case of nil weight
					if shipment.PrimeEstimatedWeight == nil {
						return nil
					}
					ihpkService, err := models.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
					suite.FatalNoError(err)

					ihpkRIOP, err := models.FetchReIntlOtherPrice(suite.DB(), *shipment.PickupAddressID, ihpkService.ID, contract.ID, shipment.RequestedPickupDate)
					suite.FatalNoError(err)
					suite.NotEmpty(ihpkRIOP)

					return models.CentPointer(computeINPKExpectedPriceCents(ihpkRIOP.PerUnitCents.Int(), escalationFactor, ntsMarketFactor, shipment.PrimeEstimatedWeight.Int()))
				}(),
			}
			suite.Equal(len(expectedServiceItems), len(serviceItems))

			// Look for INPK and assert its expected price matches the actual price the proc sets
			var foundINPK bool
			for _, serviceItem := range serviceItems {
				actualReServiceCode := serviceItem.ReService.Code
				suite.Contains(expectedServiceItems, actualReServiceCode, "Unexpected service code found: %s", actualReServiceCode)

				expectedPrice, found := expectedServiceItems[actualReServiceCode]
				suite.True(found, "Expected price for service code %s not found", actualReServiceCode)
				if actualReServiceCode == models.ReServiceCodeINPK {
					foundINPK = true
					if expectedPrice == nil || serviceItem.PricingEstimate == nil {
						// Safely ref if test case has nil outcomes
						suite.Nil(expectedPrice, "Expected price should be nil for service code %s", actualReServiceCode)
						suite.Nil(serviceItem.PricingEstimate, "Pricing estimate should be nil for service code %s", actualReServiceCode)
					} else {
						suite.Equal(*expectedPrice, *serviceItem.PricingEstimate, "Pricing estimate mismatch for service code %s", actualReServiceCode)
					}
				}
			}
			suite.FatalTrue(foundINPK)
		}
	})

}
