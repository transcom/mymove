package ppmshipment

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
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
	)

	type flags struct {
		hasSecondaryZips    bool
		hasSIT              bool
		hasProGear          bool
		hasRequestedAdvance bool
		hasReceivedAdvance  bool
	}

	// setupShipmentData - sets up old shipment based on the expected state and flags that are passed in.
	setupShipmentData := func(ppmState PPMShipmentState, oldFlags flags) (oldShipment models.PPMShipment) {
		id := uuid.Must(uuid.NewV4())
		shipmentID := uuid.Must(uuid.NewV4())

		oldShipment = models.PPMShipment{
			ID:                    id,
			ShipmentID:            shipmentID,
			Status:                models.PPMShipmentStatusDraft,
			ExpectedDepartureDate: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC),
			PickupPostalCode:      "90210",
			DestinationPostalCode: "08004",
		}

		if oldFlags.hasSecondaryZips {
			oldShipment.SecondaryPickupPostalCode = models.StringPointer("90880")
			oldShipment.SecondaryDestinationPostalCode = models.StringPointer("08900")
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
			oldShipment.ActualPickupPostalCode = models.StringPointer("90210")
			oldShipment.ActualDestinationPostalCode = models.StringPointer("79912")
			oldShipment.HasReceivedAdvance = &oldFlags.hasReceivedAdvance

			if oldFlags.hasReceivedAdvance {
				oldShipment.AdvanceAmountReceived = models.CentPointer(unit.Cents(300000))
			}
		}

		return oldShipment
	}

	type runChecksFunc func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment)

	// checkDatesAndLocationsDidntChange - ensures dates and locations fields didn't change
	checkDatesAndLocationsDidntChange := func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment) {
		suite.Equal(oldShipment.ExpectedDepartureDate, mergedShipment.ExpectedDepartureDate)
		suite.Equal(oldShipment.PickupPostalCode, mergedShipment.PickupPostalCode)
		suite.Equal(oldShipment.DestinationPostalCode, mergedShipment.DestinationPostalCode)
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
				PickupPostalCode:      "",
				DestinationPostalCode: "",
				SITExpected:           nil,
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
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
				PickupPostalCode:      "90206",
				DestinationPostalCode: "79912",
				SITExpected:           models.BoolPointer(true),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields were changed
				suite.Equal(newShipment.ExpectedDepartureDate, mergedShipment.ExpectedDepartureDate)
				suite.Equal(newShipment.PickupPostalCode, mergedShipment.PickupPostalCode)
				suite.Equal(newShipment.DestinationPostalCode, mergedShipment.DestinationPostalCode)
				suite.Equal(newShipment.SITExpected, mergedShipment.SITExpected)
			},
		},
		"Can add secondary ZIPs": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    false,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				SecondaryPickupPostalCode:      models.StringPointer("90880"),
				SecondaryDestinationPostalCode: models.StringPointer("79936"),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Equal(newShipment.SecondaryPickupPostalCode, mergedShipment.SecondaryPickupPostalCode)
				suite.Equal(newShipment.SecondaryDestinationPostalCode, mergedShipment.SecondaryDestinationPostalCode)
			},
		},
		"Can remove secondary ZIPs": {
			oldState: PPMShipmentStateDatesAndLocations,
			oldFlags: flags{
				hasSecondaryZips:    true,
				hasSIT:              false,
				hasProGear:          false,
				hasRequestedAdvance: false,
				hasReceivedAdvance:  false,
			},
			newShipment: models.PPMShipment{
				SecondaryPickupPostalCode:      models.StringPointer(""),
				SecondaryDestinationPostalCode: models.StringPointer(""),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				// ensure fields were set correctly
				suite.Nil(mergedShipment.SecondaryPickupPostalCode)
				suite.Nil(mergedShipment.SecondaryDestinationPostalCode)
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
				ActualPickupPostalCode:      models.StringPointer("90210"),
				ActualDestinationPostalCode: models.StringPointer("79912"),
				HasReceivedAdvance:          models.BoolPointer(false),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountRequested)

				// ensure fields were set correctly
				suite.Equal(newShipment.ActualPickupPostalCode, mergedShipment.ActualPickupPostalCode)
				suite.Equal(newShipment.ActualDestinationPostalCode, mergedShipment.ActualDestinationPostalCode)
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
				ActualPickupPostalCode:      models.StringPointer("90210"),
				ActualDestinationPostalCode: models.StringPointer("79912"),
				HasReceivedAdvance:          models.BoolPointer(true),
				AdvanceAmountReceived:       models.CentPointer(unit.Cents(3300)),
			},
			runChecks: func(mergedShipment models.PPMShipment, oldShipment models.PPMShipment, newShipment models.PPMShipment) {
				// ensure existing fields weren't changed
				checkDatesAndLocationsDidntChange(mergedShipment, oldShipment)
				checkEstimatedWeightsDidntChange(mergedShipment, oldShipment)
				checkSITDidntChange(mergedShipment, oldShipment)

				suite.Equal(oldShipment.HasRequestedAdvance, mergedShipment.HasRequestedAdvance)
				suite.Equal(oldShipment.AdvanceAmountRequested, mergedShipment.AdvanceAmountRequested)

				// ensure fields were set correctly
				suite.Equal(newShipment.ActualPickupPostalCode, mergedShipment.ActualPickupPostalCode)
				suite.Equal(newShipment.ActualDestinationPostalCode, mergedShipment.ActualDestinationPostalCode)
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
				suite.Equal(oldShipment.ActualPickupPostalCode, mergedShipment.ActualPickupPostalCode)
				suite.Equal(oldShipment.ActualDestinationPostalCode, mergedShipment.ActualDestinationPostalCode)

				// ensure fields were set correctly
				suite.Equal(newShipment.HasReceivedAdvance, mergedShipment.HasReceivedAdvance)
				suite.Nil(mergedShipment.AdvanceAmountReceived)
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
	}

	for name, tc := range mergeTestCases {
		name := name
		tc := tc

		suite.Run(fmt.Sprintf("Can merge changes - %s", name), func() {
			oldShipment := setupShipmentData(tc.oldState, tc.oldFlags)

			mergedShipment := mergePPMShipment(tc.newShipment, &oldShipment)

			// these should never change
			suite.Equal(oldShipment.ID, mergedShipment.ID)
			suite.Equal(oldShipment.ShipmentID, mergedShipment.ShipmentID)
			suite.Equal(oldShipment.Status, mergedShipment.Status)

			// now run test case specific checks
			tc.runChecks(*mergedShipment, oldShipment, tc.newShipment)
		})
	}
}
