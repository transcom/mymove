package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestPPMShipmentValidation() {
	validPPMShipmentStatuses := strings.Join(models.AllowedPPMShipmentStatuses, ", ")
	validPPMShipmentAdvanceStatuses := strings.Join(models.AllowedPPMAdvanceStatuses, ", ")
	validSITLocations := strings.Join(models.AllowedSITLocationTypes, ", ")

	blankAdvanceStatus := models.PPMAdvanceStatus("")
	blankSITLocation := models.SITLocationType("")
	pickupAddressID, _ := uuid.NewV4()
	destinationAddressID, _ := uuid.NewV4()

	testCases := map[string]struct {
		ppmShipment  models.PPMShipment
		expectedErrs map[string][]string
	}{
		"Successful Minimal Validation": {
			ppmShipment: models.PPMShipment{
				PPMType:               models.PPMTypeIncentiveBased,
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,
				PickupAddressID:       models.UUIDPointer(pickupAddressID),
				DestinationAddressID:  models.UUIDPointer(destinationAddressID),
			},
			expectedErrs: nil,
		},
		"Successful with optional values": {
			ppmShipment: models.PPMShipment{
				// Setting up min required fields here so that we don't get these in our errors.
				PPMType:               models.PPMTypeIncentiveBased,
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,

				// optional fields with valid values
				GunSafeWeight: models.PoundPointer(unit.Pound(500)),
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			ppmShipment: models.PPMShipment{
				PPMType:              models.PPMTypeIncentiveBased,
				PickupAddressID:      models.UUIDPointer(uuid.Nil),
				DestinationAddressID: models.UUIDPointer(uuid.Nil),
			},
			expectedErrs: map[string][]string{
				"pickup_address_id":       {"PickupAddressID can not be blank."},
				"destination_address_id":  {"DestinationAddressID can not be blank."},
				"shipment_id":             {"ShipmentID can not be blank."},
				"expected_departure_date": {"ExpectedDepartureDate can not be blank."},
				"status":                  {fmt.Sprintf("Status is not in the list [%s].", validPPMShipmentStatuses)},
			},
		},
		"Optional fields raise errors with invalid values": {
			ppmShipment: models.PPMShipment{
				// Setting up min required fields here so that we don't get these in our errors.
				PPMType:               models.PPMTypeIncentiveBased,
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,

				// Now setting optional fields with invalid values.
				DeletedAt:                 models.TimePointer(time.Time{}),
				ActualMoveDate:            models.TimePointer(time.Time{}),
				SubmittedAt:               models.TimePointer(time.Time{}),
				ReviewedAt:                models.TimePointer(time.Time{}),
				ApprovedAt:                models.TimePointer(time.Time{}),
				PickupAddressID:           models.UUIDPointer(uuid.Nil),
				DestinationAddressID:      models.UUIDPointer(uuid.Nil),
				W2AddressID:               models.UUIDPointer(uuid.Nil),
				EstimatedWeight:           models.PoundPointer(unit.Pound(-1)),
				AllowableWeight:           models.PoundPointer(unit.Pound(-1)),
				ProGearWeight:             models.PoundPointer(unit.Pound(-1)),
				SpouseProGearWeight:       models.PoundPointer(unit.Pound(-1)),
				GunSafeWeight:             models.PoundPointer(unit.Pound(-1)),
				EstimatedIncentive:        models.CentPointer(unit.Cents(-1)),
				MaxIncentive:              models.CentPointer(unit.Cents(-1)),
				FinalIncentive:            models.CentPointer(unit.Cents(0)),
				AdvanceAmountRequested:    models.CentPointer(unit.Cents(-1)),
				AdvanceStatus:             &blankAdvanceStatus,
				AdvanceAmountReceived:     models.CentPointer(unit.Cents(0)),
				SITLocation:               &blankSITLocation,
				SITEstimatedWeight:        models.PoundPointer(unit.Pound(-1)),
				SITEstimatedEntryDate:     models.TimePointer(time.Time{}),
				SITEstimatedDepartureDate: models.TimePointer(time.Time{}),
				SITEstimatedCost:          models.CentPointer(unit.Cents(0)),
				AOAPacketID:               models.UUIDPointer(uuid.Nil),
				PaymentPacketID:           models.UUIDPointer(uuid.Nil),
				GCCMultiplierID:           models.UUIDPointer(uuid.Nil),
			},
			expectedErrs: map[string][]string{
				"deleted_at":                  {"DeletedAt can not be blank."},
				"actual_move_date":            {"ActualMoveDate can not be blank."},
				"submitted_at":                {"SubmittedAt can not be blank."},
				"reviewed_at":                 {"ReviewedAt can not be blank."},
				"approved_at":                 {"ApprovedAt can not be blank."},
				"pickup_address_id":           {"PickupAddressID can not be blank."},
				"destination_address_id":      {"DestinationAddressID can not be blank."},
				"w2_address_id":               {"W2AddressID can not be blank."},
				"estimated_weight":            {"-1 is less than zero."},
				"allowable_weight":            {"-1 is less than zero."},
				"pro_gear_weight":             {"-1 is less than zero."},
				"spouse_pro_gear_weight":      {"-1 is less than zero."},
				"gun_safe_weight":             {"-1 is less than zero."},
				"estimated_incentive":         {"EstimatedIncentive cannot be negative, got: -1."},
				"max_incentive":               {"MaxIncentive cannot be negative, got: -1."},
				"final_incentive":             {"FinalIncentive must be greater than zero, got: 0."},
				"advance_amount_requested":    {"AdvanceAmountRequested cannot be negative, got: -1."},
				"advance_status":              {fmt.Sprintf("AdvanceStatus is not in the list [%s].", validPPMShipmentAdvanceStatuses)},
				"advance_amount_received":     {"AdvanceAmountReceived must be greater than zero, got: 0."},
				"sitlocation":                 {fmt.Sprintf("SITLocation is not in the list [%s].", validSITLocations)},
				"sitestimated_weight":         {"-1 is less than zero."},
				"sitestimated_entry_date":     {"SITEstimatedEntryDate can not be blank."},
				"sitestimated_departure_date": {"SITEstimatedDepartureDate can not be blank."},
				"sitestimated_cost":           {"SITEstimatedCost must be greater than zero, got: 0."},
				"aoapacket_id":                {"AOAPacketID can not be blank."},
				"payment_packet_id":           {"PaymentPacketID can not be blank."},
				"gccmultiplier_id":            {"GCCMultiplierID can not be blank."},
			},
		},
		"Gun safe weight raise error with a value over the max": {
			ppmShipment: models.PPMShipment{
				// Setting up min required fields here so that we don't get these in our errors.
				PPMType:               models.PPMTypeIncentiveBased,
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,

				// optional field with invalid value
				GunSafeWeight: models.PoundPointer(unit.Pound(501)),
			},
			expectedErrs: map[string][]string{
				"gun_safe_weight": {"must be less than or equal to 500."},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		suite.Run(name, func() {
			suite.verifyValidationErrors(testCase.ppmShipment, testCase.expectedErrs, nil)
		})
	}
}

func (suite *ModelSuite) TestCalculatePPMIncentive() {
	suite.Run("success - receive PPM incentive when all values exist", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		pickupUSPRC, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "74133", "Tulsa")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Tulsa",
					State:              "OK",
					PostalCode:         "74133",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "99505", "JBER")
		suite.FatalNoError(err)
		destAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "JBER",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		moveDate := time.Now()
		mileage := 1000
		weight := 2000

		incentives, err := models.CalculatePPMIncentive(suite.DB(), ppmShipment.ID, pickupAddress.ID, destAddress.ID, moveDate, mileage, weight, true, false, false)
		suite.NoError(err)
		suite.NotNil(incentives)
		suite.NotNil(incentives.PriceFSC)
		suite.NotNil(incentives.PriceIHPK)
		suite.NotNil(incentives.PriceIHUPK)
		suite.NotNil(incentives.PriceISLH)
		suite.NotNil(incentives.TotalIncentive)
	})

	suite.Run("failure - contract doesn't exist", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		pickupUSPRC, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "74133", "Tulsa")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Tulsa",
					State:              "OK",
					PostalCode:         "74133",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "99505", "JBER")
		suite.FatalNoError(err)
		destAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "JBER",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		// no contract for this date
		moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		mileage := 1000
		weight := 2000

		incentives, err := models.CalculatePPMIncentive(suite.DB(), ppmShipment.ID, pickupAddress.ID, destAddress.ID, moveDate, mileage, weight, true, false, false)
		suite.Error(err)
		suite.Nil(incentives)
	})
}

func (suite *ModelSuite) TestCalculatePPMSITCost() {
	suite.Run("success - receive PPM SIT costs when all values exist", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		destUSPRC, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "99505", "JBER")
		suite.FatalNoError(err)
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "JBER",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		moveDate := time.Now()
		sitDays := 7
		weight := 2000

		sitCost, err := models.CalculatePPMSITCost(suite.DB(), ppmShipment.ID, address.ID, false, moveDate, weight, sitDays)
		suite.NoError(err)
		suite.NotNil(sitCost)
		suite.NotNil(sitCost.PriceAddlDaySIT)
		suite.NotNil(sitCost.PriceFirstDaySIT)
		suite.NotNil(sitCost.TotalSITCost)
	})

	suite.Run("failure - contract doesn't exist", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		destUSPRC, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "99505", "JBER")
		suite.FatalNoError(err)
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "JBER",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		// no contract for this date
		moveDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		sitDays := 7
		weight := 2000

		sitCost, err := models.CalculatePPMSITCost(suite.DB(), ppmShipment.ID, address.ID, false, moveDate, weight, sitDays)
		suite.Error(err)
		suite.Nil(sitCost)
	})
}
