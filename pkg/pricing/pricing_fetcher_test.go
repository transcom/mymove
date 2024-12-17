package pricing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type PricingFetcherSuite struct {
	*testingsuite.PopTestSuite
}

func TestPricingFetcherSuite(t *testing.T) {

	hs := &PricingFetcherSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *PricingFetcherSuite) TestPricingFetcher() {

	setup_prices := func(isPeakPeriod bool) {

		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reServiceCodeDPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)
		reServiceCodeDDP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)
		reServiceCodeDUPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)

		startDate := time.Now().AddDate(-1, 0, 0)
		endDate := startDate.AddDate(1, 1, 1)

		contractYear := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Name:                 "Test Contract Year",
					EscalationCompounded: 1.125,
					StartDate:            startDate,
					EndDate:              endDate,
				},
			})

		serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "945",
					ServicesSchedule: 1,
				},
			})

		serviceAreaDest := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "503",
					ServicesSchedule: 1,
				},
			})

		serviceAreaPriceDOP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDOP.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(1234),
		}

		serviceAreaPriceDPK := models.ReDomesticOtherPrice{
			ContractID:   contractYear.Contract.ID,
			ServiceID:    reServiceCodeDPK.ID,
			IsPeakPeriod: isPeakPeriod,
			Schedule:     1,
			PriceCents:   unit.Cents(121),
		}

		serviceAreaPriceDDP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDDP.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaDest.ID,
			PriceCents:            unit.Cents(482),
		}

		serviceAreaPriceDUPK := models.ReDomesticOtherPrice{
			ContractID:   contractYear.Contract.ID,
			ServiceID:    reServiceCodeDUPK.ID,
			IsPeakPeriod: isPeakPeriod,
			Schedule:     1,
			PriceCents:   unit.Cents(945),
		}

		serviceAreaPriceDLH := models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			WeightLower:           500,
			WeightUpper:           10000,
			MilesLower:            1,
			MilesUpper:            10000,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceMillicents:       unit.Millicents(482),
		}

		serviceAreaPriceDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDSH.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(999),
		}

		suite.MustSave(&serviceAreaPriceDOP)
		suite.MustSave(&serviceAreaPriceDPK)
		suite.MustSave(&serviceAreaPriceDDP)
		suite.MustSave(&serviceAreaPriceDUPK)
		suite.MustSave(&serviceAreaPriceDLH)
		suite.MustSave(&serviceAreaPriceDSH)

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceArea,
				Zip3:                "945",
			},
		})

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceAreaDest,
				Zip3:                "503",
			},
		})

		testdatagen.MakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2025, time.March, 17, 0, 0, 0, 0, time.UTC),
			},
		})

	}

	suite.Run("Test Fetch Price DOP", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDOP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(69400), price)
	})

	suite.Run("Test Fetch Price DOP MTOShipmentTypeHHGOutOfNTSDom", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDOP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(76340), price)
	})

	suite.Run("Test Fetch Price DPK", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(6800), price)
	})

	suite.Run("Test Fetch Price DOP", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDDP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDDP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(27100), price)
	})

	suite.Run("Test Fetch Price DUPK", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDUPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDUPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(53150), price)
	})

	suite.Run("Test Fetch Price DLH", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDLH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			deliveryAddress.PostalCode,
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(150), price)
	})

	suite.Run("Test Fetch Price DLH No Planner", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDLH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.Error(err)
		suite.Equal(unit.Cents(0), price)
	})

	suite.Run("Test Fetch Price No estimated weight", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &setupDate,
					ShipmentType:        models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDOP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(0), price)
	})

	suite.Run("Test Fetch Price DSH", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDSH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			deliveryAddress.PostalCode,
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(281000), price)
	})

	suite.Run("Test Fetch Price DSH No Planner", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeDSH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.Error(err)
		suite.Equal(unit.Cents(0), price)
	})

	suite.Run("Test Fetch Price FSC", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeFSC,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			deliveryAddress.PostalCode,
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(7), price)
	})

	suite.Run("Test Fetch Price FSC With Actual Pickup date", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					ActualPickupDate:     &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeFSC,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			deliveryAddress.PostalCode,
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(7), price)
	})

	suite.Run("Test Fetch Price FSC No Planner", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			ReService:                 reServiceCodeFSC,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			deliveryAddress.PostalCode,
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, nil)
		// Assert

		suite.Error(err)
		suite.Equal(unit.Cents(0), price)
	})

	suite.Run("Test Fetch Price No ServiceItemCode", func() {
		// Arrange
		setup_prices(false)
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)

		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    deliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		// setup service item
		reason := "Test"

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

		mto_service_item := models.MTOServiceItem{
			MoveTaskOrder:             mto_shipment.MoveTaskOrder,
			MoveTaskOrderID:           mto_shipment.MoveTaskOrderID,
			MTOShipment:               mto_shipment,
			MTOShipmentID:             &mto_shipment.ID,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			pickupAddress.PostalCode,
			deliveryAddress.PostalCode,
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.Error(err)
		suite.Equal(unit.Cents(0), price)
	})

	suite.Run("Test LookupFSCWeightBasedDistanceMultiplier 5000", func() {
		appCtx := suite.AppContextForTest()
		weight := unit.Pound(5000)
		multiplier := LookupFSCWeightBasedDistanceMultiplier(appCtx, weight)

		suite.Equal("0.000417", multiplier)
	})
	suite.Run("Test LookupFSCWeightBasedDistanceMultiplier 10000", func() {
		appCtx := suite.AppContextForTest()
		weight := unit.Pound(10000)
		multiplier := LookupFSCWeightBasedDistanceMultiplier(appCtx, weight)

		suite.Equal("0.0006255", multiplier)
	})

	suite.Run("Test LookupFSCWeightBasedDistanceMultiplier 24000", func() {
		appCtx := suite.AppContextForTest()
		weight := unit.Pound(24000)
		multiplier := LookupFSCWeightBasedDistanceMultiplier(appCtx, weight)

		suite.Equal("0.000834", multiplier)
	})

	suite.Run("Test LookupFSCWeightBasedDistanceMultiplier 24001", func() {
		appCtx := suite.AppContextForTest()
		weight := unit.Pound(24001)
		multiplier := LookupFSCWeightBasedDistanceMultiplier(appCtx, weight)

		suite.Equal("0.00139", multiplier)
	})

	suite.Run("LookupEIAFuelPrice no value", func() {
		appCtx := suite.AppContextForTest()

		result, err := LookupEIAFuelPrice(appCtx, time.Now())

		suite.Error(err, "Looking for GHCDieselFuelPrice")
		suite.Equal(unit.Millicents(0), result)
	})

	suite.Run("Returns error when no contract year found", func() {
		setup_prices(false)
		invalidFutureDate := time.Now().AddDate(80, 0, 0)
		_, err := FetchContractCode(suite.AppContextForTest(), invalidFutureDate)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "no contract year found")
	})

	suite.Run("Returns error when no domestic service area found", func() {
		setup_prices(false)
		invalidContractCode := "invalid"
		invalidPostalCode := "00000"
		_, err := fetchDomesticServiceArea(suite.AppContextForTest(), invalidContractCode, invalidPostalCode)
		suite.Error(err)
		// suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "unable to find domestic service area for")
	})
}
