package pricing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
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

	setup_prices := func() {

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
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(1234),
		}

		serviceAreaPriceDPK := models.ReDomesticOtherPrice{
			ContractID:   contractYear.Contract.ID,
			ServiceID:    reServiceCodeDPK.ID,
			IsPeakPeriod: true,
			Schedule:     1,
			PriceCents:   unit.Cents(121),
		}

		serviceAreaPriceDDP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDDP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceAreaDest.ID,
			PriceCents:            unit.Cents(482),
		}

		serviceAreaPriceDUPK := models.ReDomesticOtherPrice{
			ContractID:   contractYear.Contract.ID,
			ServiceID:    reServiceCodeDUPK.ID,
			IsPeakPeriod: true,
			Schedule:     1,
			PriceCents:   unit.Cents(945),
		}

		serviceAreaPriceDLH := models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			WeightLower:           500,
			WeightUpper:           10000,
			MilesLower:            1,
			MilesUpper:            10000,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceMillicents:       unit.Millicents(482),
		}

		serviceAreaPriceDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDSH.ID,
			IsPeakPeriod:          true,
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
	}

	suite.Run("Test Fetch Price DOP", func() {
		// Arrange
		setup_prices()
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
		suite.Equal(unit.Cents(482), price)
	})

	suite.Run("Test Fetch Price DSH", func() {
		// Arrange
		setup_prices()
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
		suite.Equal(unit.Cents(999), price)
	})

	suite.Run("Test Fetch Price FSC", func() {
		// Arrange
		setup_prices()
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
		suite.Equal(unit.Cents(0), price)
	})

}
