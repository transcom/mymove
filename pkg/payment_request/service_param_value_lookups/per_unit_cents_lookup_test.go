package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestPerUnitCentsLookup() {
	key := models.ServiceItemParamNamePerUnitCents
	var mtoServiceItem models.MTOServiceItem
	expectedOverFiftyMilesPerUnitCentsSIT := "21796"
	expectedLessThanFiftyMilesPerUnitCentsSIT := "16417"

	setupTestData := func(serviceCode models.ReServiceCode) {
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: serviceCode,
				},
			},
		}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

	}

	setupTestDataPickupOCONUS := func(serviceCode models.ReServiceCode, sitDeliveryMileage *int) models.Move {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PickupAddressID: &address.ID,
					MarketCode:      models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		customization := make([]factory.Customization, 0)
		customization = append(customization,
			factory.Customization{
				Model:    move,
				LinkOnly: true,
			},
			factory.Customization{
				Model:    shipment,
				LinkOnly: true,
			},
			factory.Customization{
				Model: models.ReService{
					Code: serviceCode,
				},
			})

		if serviceCode == models.ReServiceCodeIOPSIT {
			customization = append(customization,
				factory.Customization{
					Model:    address,
					Type:     &factory.Addresses.SITOriginHHGActualAddress,
					LinkOnly: true,
				},
				factory.Customization{
					Model: models.MTOServiceItem{
						SITDeliveryMiles: models.IntPointer(*sitDeliveryMileage),
					},
				},
			)
		}

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), customization, nil)

		return move
	}

	setupTestDataDestOCONUS := func(serviceCode models.ReServiceCode, sitDeliveryMileage *int) models.Move {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					DestinationAddressID: &address.ID,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		customization := make([]factory.Customization, 0)
		customization = append(customization,
			factory.Customization{
				Model:    move,
				LinkOnly: true,
			},
			factory.Customization{
				Model:    shipment,
				LinkOnly: true,
			},
			factory.Customization{
				Model: models.ReService{
					Code: serviceCode,
				},
			})

		if serviceCode == models.ReServiceCodeIDDSIT {
			customization = append(customization,
				factory.Customization{
					Model:    address,
					Type:     &factory.Addresses.SITDestinationFinalAddress,
					LinkOnly: true,
				},
				factory.Customization{
					Model: models.MTOServiceItem{
						SITDeliveryMiles: models.IntPointer(*sitDeliveryMileage),
					},
				},
			)
		}

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), customization, nil)
		return move
	}

	suite.Run("success - returns perUnitCent value for IHPK", func() {
		setupTestData(models.ReServiceCodeIHPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "6997")
	})

	suite.Run("success - returns perUnitCent value for IHUPK", func() {
		setupTestData(models.ReServiceCodeIHUPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "752")
	})

	suite.Run("success - returns perUnitCent value for ISLH", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
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
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "1605")
	})

	suite.Run("success - returns perUnitCent value for IOFSIT", func() {
		move := setupTestDataPickupOCONUS(models.ReServiceCodeIOFSIT, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "607")
	})

	suite.Run("success - returns perUnitCent value for IOASIT", func() {
		move := setupTestDataPickupOCONUS(models.ReServiceCodeIOASIT, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "14")
	})

	suite.Run("success - returns perUnitCent value for IDFSIT", func() {
		move := setupTestDataDestOCONUS(models.ReServiceCodeIDFSIT, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "607")
	})

	suite.Run("success - returns perUnitCent value for IDASIT", func() {
		move := setupTestDataDestOCONUS(models.ReServiceCodeIDASIT, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "14")
	})

	suite.Run("success - returns perUnitCent value for IDASIT for a PPM", func() {
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ActualPickupDate: models.TimePointer(time.Now()),
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Tulsa",
					State:          "OK",
					PostalCode:     "74133",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "JBER",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDASIT,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCentsLookup := PerUnitCentsLookup{
			ServiceItem: mtoServiceItem,
			MTOShipment: ppm.Shipment,
		}

		appContext := suite.AppContextForTest()
		perUnitCents, err := perUnitCentsLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppm.ShipmentID,
			ContractID:    contractYear.ContractID,
		})
		suite.NoError(err)
		suite.Equal(perUnitCents, "20")
	})

	suite.Run("success - returns perUnitCent value for IUBPK", func() {
		setupTestData(models.ReServiceCodeIUBPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "7250")
	})

	suite.Run("success - returns perUnitCent value for IUBUPK", func() {
		setupTestData(models.ReServiceCodeIUBUPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "724")
	})

	suite.Run("success - returns perUnitCent value for UBP", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
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
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeUBP,
				},
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "3411")
	})

	suite.Run("success - less than 50 miles returns perUnitCent value for IOPSIT", func() {
		sitDeliveryMileage := 1
		move := setupTestDataPickupOCONUS(models.ReServiceCodeIOPSIT, models.IntPointer(sitDeliveryMileage))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(expectedLessThanFiftyMilesPerUnitCentsSIT, perUnitCents)
	})

	suite.Run("success - over 50 miles returns perUnitCent value for IOPSIT", func() {
		sitDeliveryMileage := 51
		move := setupTestDataPickupOCONUS(models.ReServiceCodeIOPSIT, models.IntPointer(sitDeliveryMileage))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(expectedOverFiftyMilesPerUnitCentsSIT, perUnitCents)
	})

	suite.Run("success - less than 50 miles returns perUnitCent value for IDDSIT", func() {
		sitDeliveryMileage := 1
		move := setupTestDataDestOCONUS(models.ReServiceCodeIDDSIT, models.IntPointer(sitDeliveryMileage))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(expectedLessThanFiftyMilesPerUnitCentsSIT, perUnitCents)
	})

	suite.Run("success - over 50 miles returns perUnitCent value for IDDSIT", func() {
		sitDeliveryMileage := 51
		move := setupTestDataDestOCONUS(models.ReServiceCodeIDDSIT, models.IntPointer(sitDeliveryMileage))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(expectedOverFiftyMilesPerUnitCentsSIT, perUnitCents)
	})

	suite.Run("failure - unauthorized service code", func() {
		setupTestData(models.ReServiceCodeDUPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Equal(perUnitCents, "")
	})

	suite.Run("failure - no requested pickup date on shipment", func() {
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHPK,
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: nil,
				},
			},
		}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

		mtoServiceItem.MTOShipment.RequestedPickupDate = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Equal(perUnitCents, "")
	})
}
