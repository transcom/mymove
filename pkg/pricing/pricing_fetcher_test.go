package pricing

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
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
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					ID:         uuid.Must(uuid.NewV4()),
					PostalCode: "945",
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					ID:         uuid.Must(uuid.NewV4()),
					PostalCode: "503",
				},
			},
		}, nil)
		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
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
		}, nil)

		// setup service item
		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reason := "Test"

		mto_service_item := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{

					MoveTaskOrderID: mto_shipment.MoveTaskOrderID,
					MTOShipmentID:   &mto_shipment.ID,
					Reason:          &reason,
					Status:          models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    mto_shipment,
				LinkOnly: true,
			},
			{
				Model:    reServiceCodeDOP,
				LinkOnly: true,
			},
		}, nil)
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
		setup_prices()
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					ID:         uuid.Must(uuid.NewV4()),
					PostalCode: "945",
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					ID:         uuid.Must(uuid.NewV4()),
					PostalCode: "503",
				},
			},
		}, nil)

		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
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
		}, nil)

		// setup service item
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)

		reason := "Test"

		mto_service_item := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					MoveTaskOrderID: mto_shipment.MoveTaskOrderID,
					MTOShipmentID:   &mto_shipment.ID,
					Reason:          &reason,
					Status:          models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    mto_shipment,
				LinkOnly: true,
			},
			{
				Model:    reServiceCodeDSH,
				LinkOnly: true,
			},
		}, nil)
		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &routemocks.Planner{}

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"945",
			"503",
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(0), price)
	})

	suite.Run("Test Fetch Price FSC", func() {
		// Arrange
		setup_prices()
		appCtx := suite.AppContextForTest()

		// setup mto shipment
		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					ID:         uuid.Must(uuid.NewV4()),
					PostalCode: "945",
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					ID:         uuid.Must(uuid.NewV4()),
					PostalCode: "503",
				},
			},
		}, nil)
		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
					RequestedPickupDate:  &setupDate,
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
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
		}, nil)

		// setup service item
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)
		reason := "Test"

		mto_service_item := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					MoveTaskOrderID: mto_shipment.MoveTaskOrderID,
					MTOShipmentID:   &mto_shipment.ID,
					Reason:          &reason,
					Status:          models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    mto_shipment,
				LinkOnly: true,
			},
			{
				Model:    reServiceCodeFSC,
				LinkOnly: true,
			},
		}, nil)
		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		suite.MustSave(&mto_shipment)

		planner := &routemocks.Planner{}

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"945",
			"503",
		).Return(5, nil)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert

		suite.NoError(err)
		suite.Equal(unit.Cents(0), price)
	})

}
