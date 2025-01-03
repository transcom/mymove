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
			"status":             {"Status is not in the list [APPROVED, REJECTED, SUBMITTED, DRAFT, CANCELLATION_REQUESTED, CANCELED, DIVERSION_REQUESTED]."},
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
	suite.Run("test MTOShipmentTypeHHGIntoNTSDom with domestic pickup and storage facility", func() {
		pickupAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		storageAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		shipment := &models.MTOShipment{
			ShipmentType:  models.MTOShipmentTypeHHGIntoNTSDom,
			PickupAddress: &pickupAddress,
			StorageFacility: &models.StorageFacility{
				Address: storageAddress,
			},
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeDomestic, updatedShipment.MarketCode, "Expected MarketCode to be d")
	})

	suite.Run("test MTOShipmentTypeHHGIntoNTSDom with international pickup", func() {
		pickupAddress := models.Address{
			IsOconus: models.BoolPointer(true),
		}
		shipment := &models.MTOShipment{
			ShipmentType:  models.MTOShipmentTypeHHGIntoNTSDom,
			PickupAddress: &pickupAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeInternational, updatedShipment.MarketCode, "Expected MarketCode to be i")
	})

	suite.Run("test MTOShipmentTypeHHGOutOfNTSDom with domestic storage and destination", func() {
		storageAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		destinationAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		shipment := &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			StorageFacility: &models.StorageFacility{
				Address: storageAddress,
			},
			DestinationAddress: &destinationAddress,
		}

		updatedShipment := models.DetermineShipmentMarketCode(shipment)
		suite.Equal(models.MarketCodeDomestic, updatedShipment.MarketCode, "Expected MarketCode to be d")
	})

	suite.Run("testMTOShipmentTypeHHGOutOfNTSDom with international destination", func() {
		storageAddress := models.Address{
			IsOconus: models.BoolPointer(false),
		}
		destinationAddress := models.Address{
			IsOconus: models.BoolPointer(true),
		}
		shipment := &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
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
