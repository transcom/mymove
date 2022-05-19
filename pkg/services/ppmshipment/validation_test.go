package ppmshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type dataSetup struct {
	old models.PPMShipment
	new models.PPMShipment
}

func setupShipmentData() (data dataSetup) {
	id := uuid.Must(uuid.NewV4())
	shipmentID := uuid.Must(uuid.NewV4())
	SITLocationOrigin := models.SITLocationTypeOrigin
	data.old = models.PPMShipment{
		ID:                    id,
		ShipmentID:            shipmentID,
		ExpectedDepartureDate: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC),
		PickupPostalCode:      "90210",
		DestinationPostalCode: "08004",
		Advance:               nil,
		AdvanceRequested:      models.BoolPointer(false),
		SitExpected:           models.BoolPointer(false),
	}
	advanceCents := unit.Cents(10000)
	estimatedWeight := unit.Pound(4000)
	proGearWeight := unit.Pound(1500)
	spouseProGearWeight := unit.Pound(400)
	data.new = models.PPMShipment{
		EstimatedWeight:           &estimatedWeight,
		HasProGear:                models.BoolPointer(true),
		ProGearWeight:             &proGearWeight,
		SpouseProGearWeight:       &spouseProGearWeight,
		Advance:                   &advanceCents,
		AdvanceRequested:          models.BoolPointer(true),
		SitExpected:               models.BoolPointer(true),
		SITEstimatedWeight:        models.PoundPointer(1000),
		SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
		SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
		SITLocation:               &SITLocationOrigin,
	}
	return data
}

func (suite *PPMShipmentSuite) TestMergePPMShipment() {
	suite.Run("Basic merge", func() {
		data := setupShipmentData()
		oldPPMShipment := data.old
		newPPMShipment := data.new

		mergedPPMShipment := mergePPMShipment(newPPMShipment, &oldPPMShipment)

		// Check if new fields are added/changed
		suite.Equal(*newPPMShipment.EstimatedWeight, *mergedPPMShipment.EstimatedWeight)
		suite.Equal(*newPPMShipment.HasProGear, *mergedPPMShipment.HasProGear)
		suite.Equal(*newPPMShipment.ProGearWeight, *mergedPPMShipment.ProGearWeight)
		suite.Equal(*newPPMShipment.SpouseProGearWeight, *mergedPPMShipment.SpouseProGearWeight)
		suite.Equal(*newPPMShipment.AdvanceRequested, *mergedPPMShipment.AdvanceRequested)
		suite.Equal(*newPPMShipment.SitExpected, *mergedPPMShipment.SitExpected)
		suite.Equal(*newPPMShipment.SITEstimatedWeight, *mergedPPMShipment.SITEstimatedWeight)
		suite.Equal(*newPPMShipment.SITEstimatedEntryDate, *mergedPPMShipment.SITEstimatedEntryDate)
		suite.Equal(*newPPMShipment.SITEstimatedDepartureDate, *mergedPPMShipment.SITEstimatedDepartureDate)
		suite.Equal(*newPPMShipment.SITLocation, *mergedPPMShipment.SITLocation)

		// Check if old fields are not changed
		suite.Equal(oldPPMShipment.ID, mergedPPMShipment.ID)
		suite.Equal(oldPPMShipment.ShipmentID, mergedPPMShipment.ShipmentID)
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, mergedPPMShipment.ExpectedDepartureDate)
		suite.Equal(oldPPMShipment.PickupPostalCode, mergedPPMShipment.PickupPostalCode)
		suite.Equal(oldPPMShipment.DestinationPostalCode, mergedPPMShipment.DestinationPostalCode)
	})

	suite.Run("Merge changes to required fields", func() {
		data := setupShipmentData()
		oldPPMShipment := data.old
		newPPMShipment := data.new
		newPPMShipment.SitExpected = models.BoolPointer(true)
		newPPMShipment.ExpectedDepartureDate = time.Date(2022, time.March, 15, 0, 0, 0, 0, time.UTC)
		newPPMShipment.PickupPostalCode = "79912"
		newPPMShipment.DestinationPostalCode = "94535"

		mergedPPMShipment := mergePPMShipment(newPPMShipment, &oldPPMShipment)

		// Check if new fields are added
		suite.Equal(newPPMShipment.ExpectedDepartureDate, mergedPPMShipment.ExpectedDepartureDate)
		suite.Equal(newPPMShipment.PickupPostalCode, mergedPPMShipment.PickupPostalCode)
		suite.Equal(newPPMShipment.DestinationPostalCode, mergedPPMShipment.DestinationPostalCode)
		suite.Equal(*newPPMShipment.SitExpected, *mergedPPMShipment.SitExpected)
		suite.Equal(*newPPMShipment.EstimatedWeight, *mergedPPMShipment.EstimatedWeight)
		suite.Equal(*newPPMShipment.HasProGear, *mergedPPMShipment.HasProGear)
		suite.Equal(*newPPMShipment.ProGearWeight, *mergedPPMShipment.ProGearWeight)
		suite.Equal(*newPPMShipment.SpouseProGearWeight, *mergedPPMShipment.SpouseProGearWeight)
		suite.Equal(*newPPMShipment.AdvanceRequested, *mergedPPMShipment.AdvanceRequested)

		// Check if old fields are not changed
		suite.Equal(oldPPMShipment.ID, mergedPPMShipment.ID)
		suite.Equal(oldPPMShipment.ShipmentID, mergedPPMShipment.ShipmentID)
	})

	suite.Run("Merge zeros", func() {
		data := setupShipmentData()
		oldPPMShipment := data.old
		oldProGearWeight := unit.Pound(1500)
		oldPPMShipment.ProGearWeight = &oldProGearWeight

		newProGearWeight := unit.Pound(0)
		newPPMShipment := models.PPMShipment{
			ProGearWeight: &newProGearWeight,
		}

		mergedPPMShipment := mergePPMShipment(newPPMShipment, &oldPPMShipment)

		// Check if new fields are updated
		suite.Equal(*newPPMShipment.ProGearWeight, *mergedPPMShipment.ProGearWeight)
	})

	suite.Run("Can remove advance", func() {
		oldPPM := models.PPMShipment{
			ID:                    uuid.Must(uuid.NewV4()),
			ShipmentID:            uuid.Must(uuid.NewV4()),
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			PickupPostalCode:      "79912",
			DestinationPostalCode: "90909",
			SitExpected:           models.BoolPointer(false),
			EstimatedWeight:       models.PoundPointer(4000),
			HasProGear:            models.BoolPointer(true),
			ProGearWeight:         models.PoundPointer(1000),
			SpouseProGearWeight:   models.PoundPointer(0),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			AdvanceRequested:      models.BoolPointer(true),
			Advance:               models.CentPointer(unit.Cents(300000)),
		}

		newPPM := models.PPMShipment{
			AdvanceRequested: models.BoolPointer(false),
		}

		mergedPPMShipment := mergePPMShipment(newPPM, &oldPPM)

		suite.Equal(*newPPM.AdvanceRequested, *mergedPPMShipment.AdvanceRequested)
		suite.Nil(mergedPPMShipment.Advance)
	})

	suite.Run("Can remove pro gear", func() {
		oldPPM := models.PPMShipment{
			ID:                    uuid.Must(uuid.NewV4()),
			ShipmentID:            uuid.Must(uuid.NewV4()),
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			PickupPostalCode:      "79912",
			DestinationPostalCode: "90909",
			SitExpected:           models.BoolPointer(false),
			EstimatedWeight:       models.PoundPointer(4000),
			HasProGear:            models.BoolPointer(true),
			ProGearWeight:         models.PoundPointer(1000),
			SpouseProGearWeight:   models.PoundPointer(0),
		}

		newPPM := models.PPMShipment{
			HasProGear: models.BoolPointer(false),
		}

		mergedPPMShipment := mergePPMShipment(newPPM, &oldPPM)

		suite.Equal(*newPPM.HasProGear, *mergedPPMShipment.HasProGear)
		suite.Nil(mergedPPMShipment.ProGearWeight)
		suite.Nil(mergedPPMShipment.SpouseProGearWeight)
	})

	suite.Run("Can remove SIT", func() {
		SITLocationOrigin := models.SITLocationTypeOrigin
		oldPPM := models.PPMShipment{
			ID:                        uuid.Must(uuid.NewV4()),
			ShipmentID:                uuid.Must(uuid.NewV4()),
			ExpectedDepartureDate:     testdatagen.NextValidMoveDate,
			PickupPostalCode:          "79912",
			DestinationPostalCode:     "90909",
			EstimatedWeight:           models.PoundPointer(4000),
			HasProGear:                models.BoolPointer(true),
			ProGearWeight:             models.PoundPointer(1000),
			SpouseProGearWeight:       models.PoundPointer(0),
			SitExpected:               models.BoolPointer(true),
			SITEstimatedWeight:        models.PoundPointer(1000),
			SITEstimatedEntryDate:     models.TimePointer(testdatagen.NextValidMoveDate),
			SITEstimatedDepartureDate: models.TimePointer(testdatagen.NextValidMoveDate.Add(testdatagen.OneWeek)),
			SITLocation:               &SITLocationOrigin,
		}

		newPPM := models.PPMShipment{
			SitExpected: models.BoolPointer(false),
		}

		mergedPPMShipment := mergePPMShipment(newPPM, &oldPPM)

		suite.Equal(*newPPM.SitExpected, *mergedPPMShipment.SitExpected)
		suite.Nil(mergedPPMShipment.SITEstimatedWeight)
		suite.Nil(mergedPPMShipment.SITEstimatedEntryDate)
		suite.Nil(mergedPPMShipment.SITEstimatedDepartureDate)
		suite.Nil(mergedPPMShipment.SITLocation)
	})

	suite.Run("Passing nil to required fields", func() {
		data := setupShipmentData()
		oldPPMShipment := data.old

		newPPMShipment := models.PPMShipment{
			ExpectedDepartureDate: time.Time{},
			PickupPostalCode:      "",
			DestinationPostalCode: "",
			SitExpected:           nil,
		}

		mergedPPMShipment := mergePPMShipment(newPPMShipment, &oldPPMShipment)

		// Check if old fields aren't updated
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, mergedPPMShipment.ExpectedDepartureDate)
		suite.Equal(oldPPMShipment.PickupPostalCode, mergedPPMShipment.PickupPostalCode)
		suite.Equal(oldPPMShipment.DestinationPostalCode, mergedPPMShipment.DestinationPostalCode)
		suite.Equal(oldPPMShipment.SitExpected, mergedPPMShipment.SitExpected)
	})
}
