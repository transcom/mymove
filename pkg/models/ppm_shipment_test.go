package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

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
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,
				PickupAddressID:       models.UUIDPointer(pickupAddressID),
				DestinationAddressID:  models.UUIDPointer(destinationAddressID),
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			ppmShipment: models.PPMShipment{
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
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,

				// Now setting optional fields with invalid values.
				DeletedAt:                   models.TimePointer(time.Time{}),
				ActualMoveDate:              models.TimePointer(time.Time{}),
				SubmittedAt:                 models.TimePointer(time.Time{}),
				ReviewedAt:                  models.TimePointer(time.Time{}),
				ApprovedAt:                  models.TimePointer(time.Time{}),
				PickupAddressID:             models.UUIDPointer(uuid.Nil),
				DestinationAddressID:        models.UUIDPointer(uuid.Nil),
				W2AddressID:                 models.UUIDPointer(uuid.Nil),
				ActualPickupPostalCode:      models.StringPointer(""),
				ActualDestinationPostalCode: models.StringPointer(""),
				EstimatedWeight:             models.PoundPointer(unit.Pound(-1)),
				ProGearWeight:               models.PoundPointer(unit.Pound(-1)),
				SpouseProGearWeight:         models.PoundPointer(unit.Pound(-1)),
				EstimatedIncentive:          models.CentPointer(unit.Cents(-1)),
				MaxIncentive:                models.CentPointer(unit.Cents(-1)),
				FinalIncentive:              models.CentPointer(unit.Cents(0)),
				AdvanceAmountRequested:      models.CentPointer(unit.Cents(-1)),
				AdvanceStatus:               &blankAdvanceStatus,
				AdvanceAmountReceived:       models.CentPointer(unit.Cents(0)),
				SITLocation:                 &blankSITLocation,
				SITEstimatedWeight:          models.PoundPointer(unit.Pound(-1)),
				SITEstimatedEntryDate:       models.TimePointer(time.Time{}),
				SITEstimatedDepartureDate:   models.TimePointer(time.Time{}),
				SITEstimatedCost:            models.CentPointer(unit.Cents(0)),
				AOAPacketID:                 models.UUIDPointer(uuid.Nil),
				PaymentPacketID:             models.UUIDPointer(uuid.Nil),
			},
			expectedErrs: map[string][]string{
				"deleted_at":                     {"DeletedAt can not be blank."},
				"actual_move_date":               {"ActualMoveDate can not be blank."},
				"submitted_at":                   {"SubmittedAt can not be blank."},
				"reviewed_at":                    {"ReviewedAt can not be blank."},
				"approved_at":                    {"ApprovedAt can not be blank."},
				"pickup_address_id":              {"PickupAddressID can not be blank."},
				"destination_address_id":         {"DestinationAddressID can not be blank."},
				"w2_address_id":                  {"W2AddressID can not be blank."},
				"actual_pickup_postal_code":      {"ActualPickupPostalCode can not be blank."},
				"actual_destination_postal_code": {"ActualDestinationPostalCode can not be blank."},
				"estimated_weight":               {"-1 is less than zero."},
				"pro_gear_weight":                {"-1 is less than zero."},
				"spouse_pro_gear_weight":         {"-1 is less than zero."},
				"estimated_incentive":            {"EstimatedIncentive cannot be negative, got: -1."},
				"max_incentive":                  {"MaxIncentive cannot be negative, got: -1."},
				"final_incentive":                {"FinalIncentive must be greater than zero, got: 0."},
				"advance_amount_requested":       {"AdvanceAmountRequested cannot be negative, got: -1."},
				"advance_status":                 {fmt.Sprintf("AdvanceStatus is not in the list [%s].", validPPMShipmentAdvanceStatuses)},
				"advance_amount_received":        {"AdvanceAmountReceived must be greater than zero, got: 0."},
				"sitlocation":                    {fmt.Sprintf("SITLocation is not in the list [%s].", validSITLocations)},
				"sitestimated_weight":            {"-1 is less than zero."},
				"sitestimated_entry_date":        {"SITEstimatedEntryDate can not be blank."},
				"sitestimated_departure_date":    {"SITEstimatedDepartureDate can not be blank."},
				"sitestimated_cost":              {"SITEstimatedCost must be greater than zero, got: 0."},
				"aoapacket_id":                   {"AOAPacketID can not be blank."},
				"payment_packet_id":              {"PaymentPacketID can not be blank."},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		suite.Run(name, func() {
			suite.verifyValidationErrors(testCase.ppmShipment, testCase.expectedErrs)
		})
	}
}
