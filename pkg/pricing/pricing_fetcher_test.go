package pricing

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FetchServiceItemPriceTestSuite) TestFetchServiceItemPrice() {

	/* 	suite.Run("Test Fetch Price Invalid Code", func() {
		// Arrange
		appCtx := suite.AppContextForTest()
		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		mto_service_item := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{},
			},
		}, nil)
		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		planner := &routemocks.Planner{}

		suite.MustSave(&mto_shipment)

		// Act
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert
		suite.Error(err)
		suite.Equal(0, price)
	}) */

	/* 	suite.Run("Test Fetch Price Bad Contract Code", func() {
		// Arrange
		appCtx := suite.AppContextForTest()
		mto_shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		planner := &routemocks.Planner{}

		suite.MustSave(&mto_shipment)

		// Act

		// mock fetch contract to fail??
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert
		suite.Error(err)
		suite.Equal(0, price)
	}) */

	suite.Run("Test Fetch Price DOP", func() {
		// Arrange
		appCtx := suite.AppContextForTest()

		setupDate := time.Now()
		estimatedWeight := unit.Pound(5000)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "23435",
				},
			},
		}, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "2343",
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

		serviceCode := models.ReServiceCodeDOP

		postalCode := "23435"
		reason := "Test"

		mto_service_item := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					MoveTaskOrder:   mto_shipment.MoveTaskOrder,
					MoveTaskOrderID: mto_shipment.MoveTaskOrderID,
					MTOShipment:     mto_shipment,
					MTOShipmentID:   &mto_shipment.ID,
					SITEntryDate:    &setupDate,
					SITPostalCode:   &postalCode,
					Reason:          &reason,
					Status:          models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    serviceCode,
				LinkOnly: true,
			},
		}, nil)
		mto_shipment.MTOServiceItems = append(mto_shipment.MTOServiceItems, mto_service_item)
		planner := &routemocks.Planner{}

		suite.MustSave(&mto_shipment)

		// Act

		// mock fetch contract to fail??
		price, err := FetchServiceItemPrice(appCtx, &mto_service_item, mto_shipment, planner)
		// Assert
		suite.Error(err)
		suite.Equal(0, price)
	})

	/*
		 	suite.Run("Test Service Item Price DOP", func() {
				// Arrange
				setupTestData := func() models.MTOShipment {
					// Set up data to use for all Origin SIT Service Item tests

					move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
					estimatedPrimeWeight := unit.Pound(6000)
					pickupDate := time.Date(2024, time.July, 31, 12, 0, 0, 0, time.UTC)
					pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
					deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

					mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
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
						{
							Model: models.MTOShipment{
								PrimeEstimatedWeight: &estimatedPrimeWeight,
								RequestedPickupDate:  &pickupDate,
							},
						},
					}, nil)

					return mtoShipment
				}

				reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
				reServiceCodeDPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)
				reServiceCodeDDP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)
				reServiceCodeDUPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)
				reServiceCodeDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
				reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)
				reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)

				startDate := time.Now().AddDate(-1, 0, 0)
				endDate := startDate.AddDate(1, 1, 1)
				sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
				sitPostalCode := "99999"
				reason := "lorem ipsum"

				contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
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

				testdatagen.MakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
					GHCDieselFuelPrice: models.GHCDieselFuelPrice{
						FuelPriceInMillicents: unit.Millicents(281400),
						PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
						EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
						EndDate:               time.Date(2025, time.March, 17, 0, 0, 0, 0, time.UTC),
					},
				})

				suite.MustSave(&serviceAreaPriceDOP)
				suite.MustSave(&serviceAreaPriceDPK)
				suite.MustSave(&serviceAreaPriceDDP)
				suite.MustSave(&serviceAreaPriceDUPK)
				suite.MustSave(&serviceAreaPriceDLH)
				suite.MustSave(&serviceAreaPriceDSH)

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

				shipment := setupTestData()
				actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
				serviceItemDOP := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeDOP,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				serviceItemDPK := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeDPK,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				serviceItemDDP := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeDDP,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				serviceItemDUPK := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeDUPK,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				serviceItemDLH := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeDLH,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				serviceItemDSH := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeDSH,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				serviceItemFSC := models.MTOServiceItem{
					MoveTaskOrder:             shipment.MoveTaskOrder,
					MoveTaskOrderID:           shipment.MoveTaskOrderID,
					MTOShipment:               shipment,
					MTOShipmentID:             &shipment.ID,
					ReService:                 reServiceCodeFSC,
					SITEntryDate:              &sitEntryDate,
					SITPostalCode:             &sitPostalCode,
					Reason:                    &reason,
					SITOriginHHGActualAddress: &actualPickupAddress,
					Status:                    models.MTOServiceItemStatusSubmitted,
				}

				testMove := setupTestData()
				estimatedSetWeight := unit.Pound(0)
				testMove.MTOServiceItems = append(testMove.MTOServiceItems, serviceItemDOP, serviceItemDDP, serviceItemDPK, serviceItemDUPK, serviceItemDLH, serviceItemDSH, serviceItemFSC)
				suite.DB().Save(testMove.MTOServiceItems)
				testMove.PrimeEstimatedWeight = &estimatedSetWeight
				suite.DB().Save(testMove)

				planner := &mocks.Planner{}

				mtoShipment := testMove.MoveTaskOrder.MTOShipments[0]

				// Setup Service item data
				// Setup Service items
				// setup mto shipment
				// Act
				// check that all service items have 0 for estimated price
				for _, serviceItem := range testMove.MTOServiceItems {
					serviceItemPrice, err := FetchServiceItemPrice(suite.AppContextForTest(), &serviceItem, mtoShipment, planner)
					suite.NoError(err)
					switch serviceItem.ReService.Code {
					case models.ReServiceCodeDOP:
						suite.Assert().Equal(unit.Cents(1234), serviceItemPrice)
					case models.ReServiceCodeDPK:
						suite.Assert().Equal(unit.Cents(121), serviceItemPrice)
					case models.ReServiceCodeDDP:
						suite.Assert().Equal(unit.Cents(482), serviceItemPrice)
					case models.ReServiceCodeDUPK:
						suite.Assert().Equal(unit.Cents(945), serviceItemPrice)
					case models.ReServiceCodeDLH:
						suite.Assert().Equal(unit.Cents(482), serviceItemPrice)
					case models.ReServiceCodeDSH:
						suite.Assert().Equal(unit.Cents(999), serviceItemPrice)
					case models.ReServiceCodeFSC:
						suite.Assert().Equal(unit.Cents(120), serviceItemPrice)
					}
				}

				// Assert

			})
	*/
}
