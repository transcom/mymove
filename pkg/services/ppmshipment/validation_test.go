package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestMergePPMShipment() {
	type PPMShipmentState int

	const (
		PPMShipmentStateDatesAndLocations         PPMShipmentState = 1
		PPMShipmentSIT                            PPMShipmentState = 2
		PPMShipmentStateEstimatedWeights          PPMShipmentState = 3
		PPMShipmentStateAdvance                   PPMShipmentState = 4
		PPMShipmentStateActualDatesZipsAndAdvance PPMShipmentState = 5
		PPMShipmentStateSecondaryAddress          PPMShipmentState = 6
		PPMShipmentStateTertiaryAddress           PPMShipmentState = 7
	)

	type flags struct {
		hasSecondaryZips               bool
		hasSIT                         bool
		hasProGear                     bool
		hasRequestedAdvance            bool
		hasReceivedAdvance             bool
		hasSecondaryPickupAddress      bool
		hasSecondaryDestinationAddress bool
		hasTertiaryPickupAddress       bool
		hasTertiaryDestinationAddress  bool
		hasActualMoveDate              bool
	}

	var (
		today      = time.Now()
		futureDate = today.AddDate(0, 0, 2)

		expectedSecondaryPickupAddress = &models.Address{
			StreetAddress1: "123 Secondary Pickup",
			City:           "New York",
			State:          "NY",
			PostalCode:     "90210",
		}
		expectedSecondaryDestinationAddress = &models.Address{
			StreetAddress1: "123 Secondary Pickup",
			City:           "New York",
			State:          "NY",
			PostalCode:     "90210",
		}
		expectedTertiaryPickupAddress = &models.Address{
			StreetAddress1: "123 Tertiary Pickup",
			City:           "New York",
			State:          "NY",
			PostalCode:     "90210",
		}
		expectedTertiaryDestinationAddress = &models.Address{
			StreetAddress1: "123 Tertiary Pickup",
			City:           "New York",
			State:          "NY",
			PostalCode:     "90210",
		}
		expectedSecondaryPickupAddressID      = uuid.Must(uuid.NewV4())
		expectedSecondaryDestinationAddressID = uuid.Must(uuid.NewV4())
		expectedTertiaryPickupAddressID       = uuid.Must(uuid.NewV4())
		expectedTertiaryDestinationAddressID  = uuid.Must(uuid.NewV4())
	)

	// setupShipmentData - sets up old shipment based on the expected state and flags that are passed in.
	setupShipmentData := func(ppmState PPMShipmentState, oldFlags flags) (oldShipment models.PPMShipment) {
		id := uuid.Must(uuid.NewV4())
		shipmentID := uuid.Must(uuid.NewV4())

		oldShipment = models.PPMShipment{
			ID:                    id,
			ShipmentID:            shipmentID,
			Status:                models.PPMShipmentStatusDraft,
			ExpectedDepartureDate: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC),
			PickupAddress: &models.Address{
				StreetAddress1: "123 Pickup",
				City:           "New York",
				State:          "NY",
				PostalCode:     "90210",
			},
			PickupAddressID: models.UUIDPointer(uuid.Must(uuid.NewV4())),
			DestinationAddress: &models.Address{
				StreetAddress1: "123 Pickup",
				City:           "New York",
				State:          "NY",
				PostalCode:     "90210",
			},
			DestinationAddressID: models.UUIDPointer(uuid.Must(uuid.NewV4())),
		}

		if ppmState >= PPMShipmentSIT {
			oldShipment.SITExpected = &oldFlags.hasSIT

			if oldFlags.hasSIT {
				SITLocationOrigin := models.SITLocationTypeOrigin
				oldShipment.SITLocation = &SITLocationOrigin
				oldShipment.SITEstimatedWeight = models.PoundPointer(unit.Pound(400))
				oldShipment.SITEstimatedEntryDate = models.TimePointer(testdatagen.NextValidMoveDate)
				oldShipment.SITEstimatedDepartureDate = models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek))
			}
		}

		if ppmState >= PPMShipmentStateEstimatedWeights {
			oldShipment.EstimatedWeight = models.PoundPointer(unit.Pound(4000))

			oldShipment.HasProGear = &oldFlags.hasProGear

			if oldFlags.hasProGear {
				oldShipment.ProGearWeight = models.PoundPointer(unit.Pound(1500))
				oldShipment.SpouseProGearWeight = models.PoundPointer(unit.Pound(400))
			}
		}

		if ppmState >= PPMShipmentStateAdvance {
			oldShipment.HasRequestedAdvance = models.BoolPointer(oldFlags.hasRequestedAdvance)

			if oldFlags.hasRequestedAdvance {
				oldShipment.AdvanceAmountRequested = models.CentPointer(unit.Cents(300000))
			}
		}

		if ppmState >= PPMShipmentStateActualDatesZipsAndAdvance {
			oldShipment.HasReceivedAdvance = &oldFlags.hasReceivedAdvance

			if oldFlags.hasReceivedAdvance {
				oldShipment.AdvanceAmountReceived = models.CentPointer(unit.Cents(300000))
			}
		}

		if ppmState >= PPMShipmentStateSecondaryAddress {
			oldShipment.HasSecondaryPickupAddress = &oldFlags.hasSecondaryPickupAddress
			if oldFlags.hasSecondaryPickupAddress {
				oldShipment.SecondaryPickupAddress = expectedSecondaryPickupAddress
				oldShipment.SecondaryPickupAddressID = &expectedSecondaryPickupAddressID
			}

			oldShipment.HasSecondaryDestinationAddress = &oldFlags.hasSecondaryDestinationAddress
			if oldFlags.hasSecondaryDestinationAddress {
				oldShipment.SecondaryDestinationAddress = expectedSecondaryDestinationAddress
				oldShipment.SecondaryDestinationAddressID = &expectedSecondaryDestinationAddressID
			}
		}

		if ppmState >= PPMShipmentStateTertiaryAddress {
			oldShipment.HasTertiaryPickupAddress = &oldFlags.hasTertiaryPickupAddress
			if oldFlags.hasTertiaryPickupAddress {
				oldShipment.TertiaryPickupAddress = expectedTertiaryPickupAddress
				oldShipment.TertiaryPickupAddressID = &expectedTertiaryPickupAddressID
			}

			oldShipment.HasTertiaryDestinationAddress = &oldFlags.hasTertiaryDestinationAddress
			if oldFlags.hasTertiaryDestinationAddress {
				oldShipment.TertiaryDestinationAddress = expectedTertiaryDestinationAddress
				oldShipment.TertiaryDestinationAddressID = &expectedTertiaryDestinationAddressID
			}
		}

		return oldShipment
	}

	type runChecksFunc func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment)

	// checkDatesAndLocationsDidntChange - ensures dates and locations fields didn't change
	checkDatesAndLocationsDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.ExpectedDepartureDate, mergedShipment.ExpectedDepartureDate)
	}

	// checkEstimatedWeightsDidntChange - ensures estimated weights fields didn't change
	checkEstimatedWeightsDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.EstimatedWeight, mergedShipment.EstimatedWeight)
		suite.Equal(oldShipment.HasProGear, mergedShipment.HasProGear)
		suite.Equal(oldShipment.ProGearWeight, mergedShipment.ProGearWeight)
		suite.Equal(oldShipment.SpouseProGearWeight, mergedShipment.SpouseProGearWeight)
	}

	// checkEstimatedWeightsDidntChange - ensures estimated weights fields didn't change
	checkSITDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.SITExpected, mergedShipment.SITExpected)
		suite.Equal(oldShipment.SITLocation, mergedShipment.SITLocation)
		suite.Equal(oldShipment.SITEstimatedWeight, mergedShipment.SITEstimatedWeight)
		suite.Equal(oldShipment.SITEstimatedEntryDate, mergedShipment.SITEstimatedEntryDate)
		suite.Equal(oldShipment.SITEstimatedDepartureDate, mergedShipment.SITEstimatedDepartureDate)
	}

	checkPickupAddressDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.PickupAddressID, mergedShipment.PickupAddressID)
		suite.Equal(oldShipment.PickupAddress.StreetAddress1, mergedShipment.PickupAddress.StreetAddress1)
		suite.Equal(oldShipment.PickupAddress.StreetAddress2, mergedShipment.PickupAddress.StreetAddress2)
		suite.Equal(oldShipment.PickupAddress.StreetAddress3, mergedShipment.PickupAddress.StreetAddress3)
		suite.Equal(oldShipment.PickupAddress.PostalCode, mergedShipment.PickupAddress.PostalCode)
		suite.Equal(oldShipment.PickupAddress.State, mergedShipment.PickupAddress.State)
	}

	checkDestinationAddressDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.DestinationAddressID, mergedShipment.DestinationAddressID)
		suite.Equal(oldShipment.DestinationAddress.StreetAddress1, mergedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(oldShipment.DestinationAddress.StreetAddress2, mergedShipment.DestinationAddress.StreetAddress2)
		suite.Equal(oldShipment.DestinationAddress.StreetAddress3, mergedShipment.DestinationAddress.StreetAddress3)
		suite.Equal(oldShipment.DestinationAddress.PostalCode, mergedShipment.DestinationAddress.PostalCode)
		suite.Equal(oldShipment.DestinationAddress.State, mergedShipment.DestinationAddress.State)
	}

	checkSecondaryPickupAddressDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.SecondaryPickupAddressID, mergedShipment.SecondaryPickupAddressID)
		suite.Equal(oldShipment.SecondaryPickupAddress.StreetAddress1, mergedShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(oldShipment.SecondaryPickupAddress.StreetAddress2, mergedShipment.SecondaryPickupAddress.StreetAddress2)
		suite.Equal(oldShipment.SecondaryPickupAddress.StreetAddress3, mergedShipment.SecondaryPickupAddress.StreetAddress3)
		suite.Equal(oldShipment.SecondaryPickupAddress.PostalCode, mergedShipment.SecondaryPickupAddress.PostalCode)
		suite.Equal(oldShipment.SecondaryPickupAddress.State, mergedShipment.SecondaryPickupAddress.State)
		suite.Equal(oldShipment.HasSecondaryPickupAddress, mergedShipment.HasSecondaryPickupAddress)
	}

	checkSecondaryDestinationAddressDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.SecondaryDestinationAddressID, mergedShipment.SecondaryDestinationAddressID)
		suite.Equal(oldShipment.SecondaryDestinationAddress.StreetAddress1, mergedShipment.SecondaryDestinationAddress.StreetAddress1)
		suite.Equal(oldShipment.SecondaryDestinationAddress.StreetAddress2, mergedShipment.SecondaryDestinationAddress.StreetAddress2)
		suite.Equal(oldShipment.SecondaryDestinationAddress.StreetAddress3, mergedShipment.SecondaryDestinationAddress.StreetAddress3)
		suite.Equal(oldShipment.SecondaryDestinationAddress.PostalCode, mergedShipment.SecondaryDestinationAddress.PostalCode)
		suite.Equal(oldShipment.SecondaryDestinationAddress.State, mergedShipment.SecondaryDestinationAddress.State)
		suite.Equal(oldShipment.SecondaryDestinationAddress.PostalCode, mergedShipment.SecondaryDestinationAddress.PostalCode)
		suite.Equal(oldShipment.HasSecondaryDestinationAddress, mergedShipment.HasSecondaryDestinationAddress)
	}

	checkTertiaryPickupAddressDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.TertiaryPickupAddressID, mergedShipment.TertiaryPickupAddressID)
		suite.Equal(oldShipment.TertiaryPickupAddress.StreetAddress1, mergedShipment.TertiaryPickupAddress.StreetAddress1)
		suite.Equal(oldShipment.TertiaryPickupAddress.StreetAddress2, mergedShipment.TertiaryPickupAddress.StreetAddress2)
		suite.Equal(oldShipment.TertiaryPickupAddress.StreetAddress3, mergedShipment.TertiaryPickupAddress.StreetAddress3)
		suite.Equal(oldShipment.TertiaryPickupAddress.PostalCode, mergedShipment.TertiaryPickupAddress.PostalCode)
		suite.Equal(oldShipment.TertiaryPickupAddress.State, mergedShipment.TertiaryPickupAddress.State)
		suite.Equal(oldShipment.HasTertiaryPickupAddress, mergedShipment.HasTertiaryPickupAddress)
	}

	checkTertiaryDestinationAddressDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.TertiaryDestinationAddressID, mergedShipment.TertiaryDestinationAddressID)
		suite.Equal(oldShipment.TertiaryDestinationAddress.StreetAddress1, mergedShipment.TertiaryDestinationAddress.StreetAddress1)
		suite.Equal(oldShipment.TertiaryDestinationAddress.StreetAddress2, mergedShipment.TertiaryDestinationAddress.StreetAddress2)
		suite.Equal(oldShipment.TertiaryDestinationAddress.StreetAddress3, mergedShipment.TertiaryDestinationAddress.StreetAddress3)
		suite.Equal(oldShipment.TertiaryDestinationAddress.PostalCode, mergedShipment.TertiaryDestinationAddress.PostalCode)
		suite.Equal(oldShipment.TertiaryDestinationAddress.State, mergedShipment.TertiaryDestinationAddress.State)
		suite.Equal(oldShipment.TertiaryDestinationAddress.PostalCode, mergedShipment.TertiaryDestinationAddress.PostalCode)
		suite.Equal(oldShipment.HasTertiaryDestinationAddress, mergedShipment.HasTertiaryDestinationAddress)
	}

	SITLocationOrigin := models.SITLocationTypeOrigin

	mergeTestCases := map[string]struct {
		oldState    PPMShipmentState
		oldFlags    flags
		newShipment models.PPMShipment
		runChecks   runChecksFunc
	}{
		"Doesn't set invalid data for required fields": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				ExpectedDepartureDate: time.Time{},
				SITExpected:           nil,
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
			},
		},
		"Edit required fields": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				ExpectedDepartureDate: time.Date(2020, time.May, 15, 0, 0, 0, 0, time.UTC),
				SITExpected:           models.BoolPointer(true),
			},
			runChecks: func(mergedShipment models.PPMShipment, _ models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields were changed
				suite.Equal(newShipment.ExpectedDepartureDate, mergedShipment.ExpectedDepartureDate)
				suite.Equal(newShipment.SITExpected, mergedShipment.SITExpected)
			},
		},
		"Add estimated weights - no pro gear": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				EstimatedWeight: models.PoundPointer(3500),
				HasProGear:      models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.EstimatedWeight, mergedShipment.EstimatedWeight)
				suite.Equal(newShipment.HasProGear, mergedShipment.HasProGear)
				suite.Nil(mergedShipment.ProGearWeight)
				suite.Nil(mergedShipment.SpouseProGearWeight)
			},
		},
		"Add estimated weights - yes pro gear": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				EstimatedWeight:     models.PoundPointer(3500),
				HasProGear:          models.BoolPointer(true),
				ProGearWeight:       models.PoundPointer(1740),
				SpouseProGearWeight: models.PoundPointer(220),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.EstimatedWeight, mergedShipment.EstimatedWeight)
				suite.Equal(newShipment.HasProGear, mergedShipment.HasProGear)
				suite.Equal(newShipment.ProGearWeight, mergedShipment.ProGearWeight)
				suite.Equal(newShipment.SpouseProGearWeight, mergedShipment.SpouseProGearWeight)
			},
		},
		"Zero out pro gear weights": {
			oldState: PPMShipmentStateEstimatedWeights,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{ // validation-wise, this shouldn't work in the end, but at the merge level, it should.
				ProGearWeight:       models.PoundPointer(unit.Pound(0)),
				SpouseProGearWeight: models.PoundPointer(unit.Pound(0)),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.EstimatedWeight, mergedShipment.EstimatedWeight)
				suite.Equal(oldShipment.HasProGear, mergedShipment.HasProGear)

				// ensure fields were set correctly
				suite.Equal(newShipment.ProGearWeight, mergedShipment.ProGearWeight)
				suite.Equal(newShipment.SpouseProGearWeight, mergedShipment.SpouseProGearWeight)
			},
		},
		"Remove pro gear": {
			oldState: PPMShipmentStateEstimatedWeights,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				HasProGear: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.EstimatedWeight, mergedShipment.EstimatedWeight)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasProGear, mergedShipment.HasProGear)
				suite.Nil(mergedShipment.ProGearWeight)
				suite.Nil(mergedShipment.SpouseProGearWeight)
			},
		},
		"Add advance requested info - no advance": {
			oldState: PPMShipmentStateEstimatedWeights,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				HasRequestedAdvance: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountRequested)
			},
		},
		"Add advance requested info - yes advance": {
			oldState: PPMShipmentStateEstimatedWeights,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				HasRequestedAdvance:    models.BoolPointer(true),
				AdvanceAmountRequested: models.CentPointer(unit.Cents(400000)),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Equal(newShipment.AdvanceAmountRequested, mergedShipment.AdvanceAmountRequested)
			},
		},
		"Remove advance requested": {
			oldState: PPMShipmentStateAdvance,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: true,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				HasRequestedAdvance: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountRequested)
			},
		},
		"Add actual zips and advance info - no advance": {
			oldState: PPMShipmentStateAdvance,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				HasReceivedAdvance: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountRequested)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasReceivedAdvance, mergedShipment.HasReceivedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountReceived)
			},
		},
		"Add actual zips and advance info - yes advance": {
			oldState: PPMShipmentStateAdvance,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: true,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				HasReceivedAdvance:    models.BoolPointer(true),
				AdvanceAmountReceived: models.CentPointer(unit.Cents(3300)),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Equal(oldShipment.AdvanceAmountRequested, mergedShipment.AdvanceAmountRequested)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasReceivedAdvance, mergedShipment.HasReceivedAdvance)
				suite.Equal(newShipment.AdvanceAmountReceived, mergedShipment.AdvanceAmountReceived)
			},
		},
		"Remove actual advance": {
			oldState: PPMShipmentStateActualDatesZipsAndAdvance,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: true,
				hasReceivedAdvance:  true,
			},
			newShipment: models.PPMShipment{
				HasReceivedAdvance: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Equal(oldShipment.AdvanceAmountRequested, mergedShipment.AdvanceAmountRequested)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasReceivedAdvance, mergedShipment.HasReceivedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountReceived)
			},
		},
		"Add W2 Address and Final Incentive": {
			oldState: PPMShipmentStateActualDatesZipsAndAdvance,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          true,
				hasRequestedAdvance: true,
				hasReceivedAdvance:  true,
			},
			newShipment: models.PPMShipment{
				W2Address: &models.Address{
					StreetAddress1: "123 Main",
					City:           "New York",
					State:          "NY",
					PostalCode:     "90210",
				},
				FinalIncentive: models.CentPointer(unit.Cents(3300)),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.FinalIncentive, mergedShipment.FinalIncentive)
				suite.Equal(newShipment.W2Address.StreetAddress1, mergedShipment.W2Address.StreetAddress1)
				suite.Equal(newShipment.W2Address.City, mergedShipment.W2Address.City)
				suite.Equal(newShipment.W2Address.PostalCode, mergedShipment.W2Address.PostalCode)
				suite.Equal(newShipment.W2Address.State, mergedShipment.W2Address.State)
			},
		},
		"Add SIT info - SIT expected ": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				SITExpected:               models.BoolPointer(true),
				SITLocation:               &SITLocationOrigin,
				SITEstimatedWeight:        models.PoundPointer(unit.Pound(400)),
				SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
				SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.SITExpected, mergedShipment.SITExpected)
				suite.Equal(newShipment.SITLocation, mergedShipment.SITLocation)
				suite.Equal(newShipment.SITEstimatedWeight, mergedShipment.SITEstimatedWeight)
				suite.Equal(newShipment.SITEstimatedEntryDate, mergedShipment.SITEstimatedEntryDate)
				suite.Equal(newShipment.SITEstimatedDepartureDate, mergedShipment.SITEstimatedDepartureDate)
			},
		},
		"Remove SIT info - SIT not expected": {
			oldState: PPMShipmentStateEstimatedWeights,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              true,
				hasProGear:          true,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				SITExpected: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.SITExpected, mergedShipment.SITExpected)
				suite.Nil(mergedShipment.SITLocation)
				suite.Nil(mergedShipment.SITEstimatedWeight)
				suite.Nil(mergedShipment.SITEstimatedEntryDate)
				suite.Nil(mergedShipment.SITEstimatedDepartureDate)
			},
		},
		"default HasSecondaryPickupAddress and HasSecondaryDestinationAddress": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				//hasSecondaryPickupAddress and hasSecondaryDestinationAddress not provided, assumes no deletes
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryPickupAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)
			},
		},
		"default HasTertiaryPickupAddress and HasTertiaryDestinationAddress": {
			oldState: PPMShipmentStateTertiaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
				hasTertiaryPickupAddress:       true,
				hasTertiaryDestinationAddress:  true,
			},
			newShipment: models.PPMShipment{
				//hasTertiaryPickupAddress and hasTertiaryDestinationAddress not provided, assumes no deletes
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryPickupAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkTertiaryPickupAddressDidntChange(mergedShipment, oldShipment)
				checkTertiaryDestinationAddressDidntChange(mergedShipment, oldShipment)
			},
		},
		"delete secondaryPickAddress/ID by HasSecondaryPickupAddress": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				HasSecondaryPickupAddress: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)

				// True flag will null out both ID and model
				suite.True(mergedShipment.SecondaryPickupAddress == nil)
				suite.True(mergedShipment.SecondaryPickupAddressID == nil)
			},
		},
		"delete tertiaryPickAddress/ID by HasTertiaryPickupAddress": {
			oldState: PPMShipmentStateTertiaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
				hasTertiaryPickupAddress:       true,
				hasTertiaryDestinationAddress:  true,
			},
			newShipment: models.PPMShipment{
				HasTertiaryPickupAddress: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryPickupAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkTertiaryDestinationAddressDidntChange(mergedShipment, oldShipment)

				// True flag will null out both ID and model
				suite.True(mergedShipment.TertiaryPickupAddress == nil)
				suite.True(mergedShipment.TertiaryPickupAddressID == nil)
			},
		},
		"delete secondaryDestinationAddress/ID by HasSecondaryDestinationAddress": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				HasSecondaryDestinationAddress: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				// ensure existing fields weren't changed
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryPickupAddressDidntChange(mergedShipment, oldShipment)

				// True flag will null out both ID and model
				suite.True(mergedShipment.SecondaryDestinationAddress == nil)
				suite.True(mergedShipment.SecondaryDestinationAddressID == nil)
			},
		},
		"delete tertiaryDestinationAddress/ID by HasTertiaryDestinationAddress": {
			oldState: PPMShipmentStateTertiaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
				hasTertiaryPickupAddress:       true,
				hasTertiaryDestinationAddress:  true,
			},
			newShipment: models.PPMShipment{
				HasTertiaryDestinationAddress: models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkTertiaryPickupAddressDidntChange(mergedShipment, oldShipment)

				// True flag will null out both ID and model
				suite.True(mergedShipment.TertiaryDestinationAddress == nil)
				suite.True(mergedShipment.TertiaryDestinationAddressID == nil)
			},
		},
		"update secondaryPickupAddress with no HasSecondaryPickupAddress=true": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				SecondaryPickupAddress: &models.Address{
					StreetAddress1: "updated",
					City:           "updated",
					State:          "NY",
					PostalCode:     "11111",
				},
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.SecondaryPickupAddress.City, mergedShipment.SecondaryPickupAddress.City)
				suite.Equal(newShipment.SecondaryPickupAddress.StreetAddress1, mergedShipment.SecondaryPickupAddress.StreetAddress1)
				suite.Equal(newShipment.SecondaryPickupAddress.PostalCode, mergedShipment.SecondaryPickupAddress.PostalCode)
				suite.Equal(oldShipment.HasSecondaryPickupAddress, mergedShipment.HasSecondaryPickupAddress)
				suite.Equal(oldShipment.SecondaryPickupAddressID, mergedShipment.SecondaryPickupAddressID)
			},
		},
		"update tertiaryPickupAddress with no HasTertiaryPickupAddress=true": {
			oldState: PPMShipmentStateTertiaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
				hasTertiaryPickupAddress:       true,
				hasTertiaryDestinationAddress:  true,
			},
			newShipment: models.PPMShipment{
				TertiaryPickupAddress: &models.Address{
					StreetAddress1: "updated",
					City:           "updated",
					State:          "NY",
					PostalCode:     "11111",
				},
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryPickupAddressDidntChange(mergedShipment, oldShipment)
				checkTertiaryDestinationAddressDidntChange(mergedShipment, oldShipment)
				// ensure fields were set correctly
				suite.Equal(newShipment.TertiaryPickupAddress.City, mergedShipment.TertiaryPickupAddress.City)
				suite.Equal(newShipment.TertiaryPickupAddress.StreetAddress1, mergedShipment.TertiaryPickupAddress.StreetAddress1)
				suite.Equal(newShipment.TertiaryPickupAddress.PostalCode, mergedShipment.TertiaryPickupAddress.PostalCode)
				suite.Equal(oldShipment.HasTertiaryPickupAddress, mergedShipment.HasTertiaryPickupAddress)
				suite.Equal(oldShipment.TertiaryPickupAddressID, mergedShipment.TertiaryPickupAddressID)
			},
		},
		"update secondaryPickupAddress with HasSecondaryPickupAddress=true": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				HasSecondaryPickupAddress: models.BoolPointer(true),
				SecondaryPickupAddress: &models.Address{
					StreetAddress1: "updated",
					City:           "updated",
					State:          "NY",
					PostalCode:     "11111",
				},
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.SecondaryPickupAddress.City, mergedShipment.SecondaryPickupAddress.City)
				suite.Equal(newShipment.SecondaryPickupAddress.StreetAddress1, mergedShipment.SecondaryPickupAddress.StreetAddress1)
				suite.Equal(newShipment.SecondaryPickupAddress.PostalCode, mergedShipment.SecondaryPickupAddress.PostalCode)
				suite.Equal(oldShipment.HasSecondaryPickupAddress, mergedShipment.HasSecondaryPickupAddress)
				suite.Equal(oldShipment.SecondaryPickupAddressID, mergedShipment.SecondaryPickupAddressID)

			},
		},
		"attempt to update secondaryPickupAddress with HasSecondaryPickupAddress=false": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				HasSecondaryPickupAddress: models.BoolPointer(false),
				// this should be ignored
				SecondaryPickupAddress: &models.Address{
					StreetAddress1: "updated",
					City:           "updated",
					State:          "NY",
					PostalCode:     "11111",
				},
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				// ensure existing fields weren't changed
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryDestinationAddressDidntChange(mergedShipment, oldShipment)

				// verify delete occured
				suite.True(mergedShipment.SecondaryPickupAddress == nil)
				suite.True(mergedShipment.SecondaryPickupAddressID == nil)
			},
		},
		"attempt to update secondaryDestinationAddress with HasSecondaryDestinationAddress=false": {
			oldState: PPMShipmentStateSecondaryAddress,
			oldFlags: flags{
				hasSecondaryPickupAddress:      true,
				hasSecondaryDestinationAddress: true,
			},
			newShipment: models.PPMShipment{
				HasSecondaryDestinationAddress: models.BoolPointer(false),
				// this should be ignored
				SecondaryDestinationAddress: &models.Address{
					StreetAddress1: "updated",
					City:           "updated",
					State:          "NY",
					PostalCode:     "11111",
				},
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, _ models.PPMShipment) {
				// ensure existing fields weren't changed
				checkPickupAddressDidntChange(mergedShipment, oldShipment)
				checkDestinationAddressDidntChange(mergedShipment, oldShipment)
				checkSecondaryPickupAddressDidntChange(mergedShipment, oldShipment)

				// verify delete occured
				suite.True(mergedShipment.SecondaryDestinationAddress == nil)
				suite.True(mergedShipment.SecondaryDestinationAddressID == nil)
			},
		},
		"attempt to update actual move date with invalid date": {
			oldFlags: flags{
				hasActualMoveDate: true,
			},
			newShipment: models.PPMShipment{
				ActualMoveDate: &futureDate,
			},
			runChecks: func(_ models.PPMShipment, _ models.PPMShipment, _ models.PPMShipment) {
			},
		},
	}

	for name, tc := range mergeTestCases {
		name := name
		tc := tc

		suite.Run(fmt.Sprintf("Can merge changes - %s", name), func() {
			oldShipment := setupShipmentData(tc.oldState, tc.oldFlags)

			mergedShipment, err := mergePPMShipment(tc.newShipment, &oldShipment)

			// these should never change
			suite.Equal(oldShipment.ID, mergedShipment.ID)
			suite.Equal(oldShipment.ShipmentID, mergedShipment.ShipmentID)
			suite.Equal(oldShipment.Status, mergedShipment.Status)

			if tc.oldFlags.hasActualMoveDate {
				suite.Equal(err.Error(), "Update Error Actual move date cannot be set to the future.")
			}

			// now run test case specific checks
			tc.runChecks(*mergedShipment, oldShipment, tc.newShipment)
		})
	}
}
