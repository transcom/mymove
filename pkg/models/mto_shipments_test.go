package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestMTOShipmentValidation() {
	suite.Run("test valid MTOShipment", func() {
		// mock weights
		estimatedWeight := unit.Pound(1000)
		actualWeight := unit.Pound(980)
		sitDaysAllowance := 90
		tacType := models.LOATypeHHG
		sacType := models.LOATypeHHG
		marketCode := models.MarketCodeDomestic
		validMTOShipment := models.MTOShipment{
			MoveTaskOrderID:      uuid.Must(uuid.NewV4()),
			Status:               models.MTOShipmentStatusApproved,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			SITDaysAllowance:     &sitDaysAllowance,
			TACType:              &tacType,
			SACType:              &sacType,
			MarketCode:           marketCode,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOShipment, expErrors)
	})

	suite.Run("test empty MTOShipment", func() {
		emptyMTOShipment := models.MTOShipment{}
		expErrors := map[string][]string{
			"move_task_order_id": {"MoveTaskOrderID can not be blank."},
			"status":             {"Status is not in the list [APPROVED, REJECTED, SUBMITTED, DRAFT, CANCELLATION_REQUESTED, CANCELED, DIVERSION_REQUESTED, TERMINATED_FOR_CAUSE, APPROVALS_REQUESTED]."},
		}
		suite.verifyValidationErrors(&emptyMTOShipment, expErrors)
	})

	suite.Run("test rejected MTOShipment", func() {
		rejectionReason := "bad shipment"
		marketCode := models.MarketCodeDomestic
		rejectedMTOShipment := models.MTOShipment{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			Status:          models.MTOShipmentStatusRejected,
			MarketCode:      marketCode,
			RejectionReason: &rejectionReason,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&rejectedMTOShipment, expErrors)
	})

	suite.Run("test validation failures", func() {
		// mock weights
		estimatedWeight := unit.Pound(-1000)
		actualWeight := unit.Pound(-980)
		billableWeightCap := unit.Pound(-1)
		billableWeightJustification := ""
		sitDaysAllowance := -1
		serviceOrderNumber := ""
		tacType := models.LOAType("FAKE")
		marketCode := models.MarketCode("x")
		invalidMTOShipment := models.MTOShipment{
			MoveTaskOrderID:             uuid.Must(uuid.NewV4()),
			Status:                      models.MTOShipmentStatusRejected,
			PrimeEstimatedWeight:        &estimatedWeight,
			PrimeActualWeight:           &actualWeight,
			BillableWeightCap:           &billableWeightCap,
			BillableWeightJustification: &billableWeightJustification,
			SITDaysAllowance:            &sitDaysAllowance,
			ServiceOrderNumber:          &serviceOrderNumber,
			StorageFacilityID:           &uuid.Nil,
			TACType:                     &tacType,
			SACType:                     &tacType,
			MarketCode:                  marketCode,
		}
		expErrors := map[string][]string{
			"prime_estimated_weight":        {"-1000 is not greater than 0."},
			"prime_actual_weight":           {"-980 is not greater than 0."},
			"rejection_reason":              {"RejectionReason can not be blank."},
			"billable_weight_cap":           {"-1 is less than zero."},
			"billable_weight_justification": {"BillableWeightJustification can not be blank."},
			"sitdays_allowance":             {"-1 is not greater than -1."},
			"service_order_number":          {"ServiceOrderNumber can not be blank."},
			"storage_facility_id":           {"StorageFacilityID can not be blank."},
			"tactype":                       {"TACType is not in the list [HHG, NTS]."},
			"sactype":                       {"SACType is not in the list [HHG, NTS]."},
			"market_code":                   {"MarketCode is not in the list [d, i]."},
		}
		suite.verifyValidationErrors(&invalidMTOShipment, expErrors)
	})
	suite.Run("test MTO Shipment has a PPM Shipment", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		mtoShipment.PPMShipment = &ppmShipment
		result := mtoShipment.ContainsAPPMShipment()

		suite.True(result, "Expected mtoShipment to cotain a PPM Shipment")
	})
}

func (suite *ModelSuite) TestDetermineShipmentMarketCode() {
	suite.Run("test MTOShipmentTypeHHGIntoNTS with domestic pickup and storage facility", func() {
		pickupAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		storageAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		shipment := &models.MTOShipment{
			ShipmentType:  models.MTOShipmentTypeHHGIntoNTS,
			PickupAddress: &pickupAddress,
			StorageFacility: &models.StorageFacility{
				Address: storageAddress,
			},
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeDomestic, updatedShipment.MarketCode, "Expected MarketCode to be d")
	})

	suite.Run("test MTOShipmentTypeHHGIntoNTS with international pickup", func() {
		pickupAddress := models.Address{
			IsOconus: models.BoolPointer(true),
		}
		shipment := &models.MTOShipment{
			ShipmentType:  models.MTOShipmentTypeHHGIntoNTS,
			PickupAddress: &pickupAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeInternational, updatedShipment.MarketCode, "Expected MarketCode to be i")
	})

	suite.Run("test MTOShipmentTypeHHGOutOfNTS with domestic storage and destination", func() {
		storageAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		destinationAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		shipment := &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGOutOfNTS,
			StorageFacility: &models.StorageFacility{
				Address: storageAddress,
			},
			DestinationAddress: &destinationAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeDomestic, updatedShipment.MarketCode, "Expected MarketCode to be d")
	})

	suite.Run("testMTOShipmentTypeHHGOutOfNTS with international destination", func() {
		storageAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		destinationAddress := models.Address{
			IsOconus: models.BoolPointer(true),
		}
		shipment := &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGOutOfNTS,
			StorageFacility: &models.StorageFacility{
				Address: storageAddress,
			},
			DestinationAddress: &destinationAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeInternational, updatedShipment.MarketCode, "Expected MarketCode to be i")
	})

	suite.Run("test default shipment with domestic pickup and destination", func() {
		pickupAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		destinationAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		shipment := &models.MTOShipment{
			PickupAddress:      &pickupAddress,
			DestinationAddress: &destinationAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeDomestic, updatedShipment.MarketCode, "Expected MarketCode to be d")
	})

	suite.Run("test default shipment with international destination", func() {
		pickupAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		destinationAddress := models.Address{
			IsOconus: models.BoolPointer(true),
		}
		shipment := &models.MTOShipment{
			PickupAddress:      &pickupAddress,
			DestinationAddress: &destinationAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeInternational, updatedShipment.MarketCode, "Expected MarketCode to be i")
	})
}

func (suite *ModelSuite) TestDetermineMarketCode() {
	marketCodeNil := models.MarketCode("")
	suite.Run("test domestic market code for two CONUS addresses", func() {
		address1 := &models.Address{
			IsOconus: models.BoolPointer(false),
		}
		address2 := &models.Address{
			IsOconus: models.BoolPointer(false),
		}

		marketCode, err := models.DetermineMarketCode(address1, address2)
		suite.NoError(err)
		suite.Equal(models.MarketCodeDomestic, marketCode, "Expected MarketCode to be d")
	})

	suite.Run("test international market code with CONUS and OCONUS address", func() {
		address1 := &models.Address{
			IsOconus: models.BoolPointer(false),
		}
		address2 := &models.Address{
			IsOconus: models.BoolPointer(true),
		}

		marketCode, err := models.DetermineMarketCode(address1, address2)
		suite.NoError(err)
		suite.Equal(models.MarketCodeInternational, marketCode, "Expected MarketCode to be i")
	})

	suite.Run("test international market code for two OCONUS addresses", func() {
		address1 := &models.Address{
			IsOconus: models.BoolPointer(true),
		}
		address2 := &models.Address{
			IsOconus: models.BoolPointer(true),
		}

		marketCode, err := models.DetermineMarketCode(address1, address2)
		suite.NoError(err)
		suite.Equal(models.MarketCodeInternational, marketCode, "Expected MarketCode to be i")
	})

	suite.Run("test error when address1 is nil", func() {
		address1 := (*models.Address)(nil)
		address2 := &models.Address{
			IsOconus: models.BoolPointer(false),
		}

		marketCode, err := models.DetermineMarketCode(address1, address2)
		suite.Error(err)
		suite.Equal(marketCodeNil, marketCode, "Expected MarketCode to be empty when address1 is nil")
		suite.EqualError(err, "both address1 and address2 must be provided")
	})

	suite.Run("test error when address2 is nil", func() {
		address1 := &models.Address{
			IsOconus: models.BoolPointer(false),
		}
		address2 := (*models.Address)(nil)

		marketCode, err := models.DetermineMarketCode(address1, address2)
		suite.Error(err)
		suite.Equal(marketCodeNil, marketCode, "Expected MarketCode to be empty when address2 is nil")
		suite.EqualError(err, "both address1 and address2 must be provided")
	})

	suite.Run("test error when both addresses are nil", func() {
		address1 := (*models.Address)(nil)
		address2 := (*models.Address)(nil)

		marketCode, err := models.DetermineMarketCode(address1, address2)
		suite.Error(err)
		suite.Equal(marketCodeNil, marketCode, "Expected MarketCode to be empty when both addresses are nil")
		suite.EqualError(err, "both address1 and address2 must be provided")
	})
}
func (suite *ModelSuite) TestCreateApprovedServiceItemsForShipment() {
	suite.Run("test creating approved service items for shipment", func() {

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{

				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "CA",
					PostalCode:     "90210",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.MTOShipment{
					MarketCode: "i",
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)
		err := models.CreateApprovedServiceItemsForShipment(suite.DB(), &shipment)
		suite.NoError(err)
	})

	suite.Run("test error handling for invalid shipment", func() {
		invalidShipment := models.MTOShipment{}

		err := models.CreateApprovedServiceItemsForShipment(suite.DB(), &invalidShipment)
		suite.Error(err)
	})
}

func (suite *ModelSuite) TestCreateInternationalAccessorialServiceItemsForShipment() {
	suite.Run("test creating accessorial service items for shipment", func() {

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					RejectionReason: models.StringPointer("not applicable"),
					MTOShipmentID:   &shipment.ID,
					Reason:          models.StringPointer("this is a special item"),
					EstimatedWeight: models.PoundPointer(400),
					ActualWeight:    models.PoundPointer(500),
				},
			},
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
					Code: models.ReServiceCodeIDSHUT,
				},
			},
		}, nil)

		serviceItem.MTOShipment = shipment
		serviceItemIds, err := models.CreateInternationalAccessorialServiceItemsForShipment(suite.DB(), shipment.ID, models.MTOServiceItems{serviceItem})
		suite.NoError(err)
		suite.NotNil(serviceItemIds)
	})

	suite.Run("test error handling for invalid shipment", func() {
		serviceItemIds, err := models.CreateInternationalAccessorialServiceItemsForShipment(suite.DB(), uuid.Nil, models.MTOServiceItems{})
		suite.Error(err)
		suite.Nil(serviceItemIds)
	})
}

func (suite *ModelSuite) TestFindShipmentByID() {
	suite.Run("success - test find", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		_, err := models.FetchShipmentByID(suite.DB(), shipment.ID)
		suite.NoError(err)
	})

	suite.Run("not found test find", func() {
		notValidID := uuid.Must(uuid.NewV4())
		_, err := models.FetchShipmentByID(suite.DB(), notValidID)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound, err)
	})
}

func (suite *ModelSuite) TestGetDestinationGblocForShipment() {
	suite.Run("success - get GBLOC for USAF in AK Zone II", func() {
		// Create a USAF move in Alaska Zone II
		// this is a hard coded uuid that is a us_post_region_cities_id within AK Zone II
		// this should always return MBFL
		zone2UUID, err := uuid.FromString("66768964-e0de-41f3-b9be-7ef32e4ae2b4")
		suite.FatalNoError(err)
		airForce := models.AffiliationAIRFORCE
		postalCode := "99501"

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode:         postalCode,
					UsPostRegionCityID: &zone2UUID,
				},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &airForce,
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		gbloc, err := models.GetDestinationGblocForShipment(suite.DB(), shipment.ID)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(*gbloc, "MBFL")
	})
	suite.Run("success - get GBLOC for Army in AK Zone II", func() {
		// Create an ARMY move in Alaska Zone II
		zone2UUID, err := uuid.FromString("66768964-e0de-41f3-b9be-7ef32e4ae2b4")
		suite.FatalNoError(err)
		army := models.AffiliationARMY
		postalCode := "99501"
		// since we truncate the test db, we need to add the postal_code_to_gbloc value
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "99744", "JEAT")

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode:         postalCode,
					UsPostRegionCityID: &zone2UUID,
				},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		gbloc, err := models.GetDestinationGblocForShipment(suite.DB(), shipment.ID)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(*gbloc, "JEAT")
	})
	suite.Run("success - get GBLOC for USMC in AK Zone II", func() {
		// Create a USMC move in Alaska Zone II
		// this should always return USMC
		zone2UUID, err := uuid.FromString("66768964-e0de-41f3-b9be-7ef32e4ae2b4")
		suite.FatalNoError(err)
		usmc := models.AffiliationMARINES
		postalCode := "99501"
		// since we truncate the test db, we need to add the postal_code_to_gbloc value
		// this doesn't matter to the db function because it will check for USMC but we are just verifying it won't be JEAT despite the zip matching
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "99744", "JEAT")

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode:         postalCode,
					UsPostRegionCityID: &zone2UUID,
				},
			},
		}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &usmc,
				},
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
			},
		}, nil)

		gbloc, err := models.GetDestinationGblocForShipment(suite.DB(), shipment.ID)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(*gbloc, "USMC")
	})
}

func (suite *ModelSuite) TestIsPPMShipment() {
	suite.Run("true - shipment is a ppm", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		mtoShipment.PPMShipment = &ppmShipment
		mtoShipment.ShipmentType = models.MTOShipmentTypePPM

		isPPM := mtoShipment.IsPPMShipment()
		suite.NotNil(isPPM)
		suite.Equal(isPPM, true)
	})

	suite.Run("false - shipment is not a ppm", func() {
		nonPPMshipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		isPPM := nonPPMshipment.IsPPMShipment()
		suite.NotNil(isPPM)
		suite.Equal(isPPM, false)
	})
}

func (suite *ModelSuite) TestIsShipmentOCONUS() {
	suite.Run("dest OCONUS but pickup CONUS", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "CA",
					PostalCode:     "90210",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		isOCONUS := models.IsShipmentOCONUS(shipment)
		suite.NotNil(isOCONUS)
		suite.True(*isOCONUS)
	})

	suite.Run("pickup OCONUS but dest CONUS", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "CA",
					PostalCode:     "90210",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.PickupAddress,
			},
		}, nil)

		isOCONUS := models.IsShipmentOCONUS(shipment)
		suite.NotNil(isOCONUS)
		suite.True(*isOCONUS)
	})

	suite.Run("pickup CONUS, dest CONUS", func() {
		// default factory produces two CONUS addresses
		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		isOCONUS := models.IsShipmentOCONUS(shipment)
		suite.NotNil(isOCONUS)
		suite.False(*isOCONUS)
	})

	suite.Run("both OCONUS addresses", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some other address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.PickupAddress,
			},
		}, nil)

		isOCONUS := models.IsShipmentOCONUS(shipment)
		suite.NotNil(isOCONUS)
		suite.True(*isOCONUS)
	})

	suite.Run("nil PickupAddress.IsOconus", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "CA",
					PostalCode:     "90210",
					IsOconus:       nil,
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		shipment.PickupAddress.IsOconus = nil

		isOCONUS := models.IsShipmentOCONUS(shipment)
		suite.Nil(isOCONUS)
	})

	suite.Run("nil DestinationAddress.IsOconus", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "CA",
					PostalCode:     "90210",
					IsOconus:       models.BoolPointer(false),
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "some address",
					City:           "city",
					State:          "AK",
					PostalCode:     "98765",
					IsOconus:       nil,
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		shipment.DestinationAddress.IsOconus = nil

		isOCONUS := models.IsShipmentOCONUS(shipment)
		suite.Nil(isOCONUS)
	})
}
