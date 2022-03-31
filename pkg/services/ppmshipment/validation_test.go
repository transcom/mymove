package ppmshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type dataSetup struct {
	old models.PPMShipment
	new models.PPMShipment
}

func setupShipmentData() (data dataSetup) {
	id := uuid.Must(uuid.NewV4())
	shipmentID := uuid.Must(uuid.NewV4())
	data.old = models.PPMShipment{
		ID:                    id,
		ShipmentID:            shipmentID,
		ExpectedDepartureDate: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC),
		PickupPostalCode:      "90210",
		DestinationPostalCode: "08004",
		SitExpected:           models.BoolPointer(false),
		Advance:               nil,
		AdvanceRequested:      models.BoolPointer(false),
	}
	advanceCents := unit.Cents(10000)
	estimatedWeight := unit.Pound(4000)
	proGearWeight := unit.Pound(1500)
	spouseProGearWeight := unit.Pound(400)
	data.new = models.PPMShipment{
		EstimatedWeight:     &estimatedWeight,
		HasProGear:          models.BoolPointer(true),
		ProGearWeight:       &proGearWeight,
		SpouseProGearWeight: &spouseProGearWeight,
		Advance:             &advanceCents,
		AdvanceRequested:    models.BoolPointer(true),
	}
	return data
}

func (suite *PPMShipmentSuite) TestMergePPMShipment() {
	suite.Run("Basic merge", func() {
		data := setupShipmentData()
		oldPPMShipment := data.old
		newPPMShipment := data.new

		mergedPPMShipment := mergePPMShipment(newPPMShipment, &oldPPMShipment)

		// Check if new fields are added
		suite.Equal(*newPPMShipment.EstimatedWeight, *mergedPPMShipment.EstimatedWeight)
		suite.Equal(*newPPMShipment.HasProGear, *mergedPPMShipment.HasProGear)
		suite.Equal(*newPPMShipment.ProGearWeight, *mergedPPMShipment.ProGearWeight)
		suite.Equal(*newPPMShipment.SpouseProGearWeight, *mergedPPMShipment.SpouseProGearWeight)
		suite.Equal(*newPPMShipment.AdvanceRequested, *mergedPPMShipment.AdvanceRequested)

		// Check if old fields are not changed
		suite.Equal(oldPPMShipment.ID, mergedPPMShipment.ID)
		suite.Equal(oldPPMShipment.ShipmentID, mergedPPMShipment.ShipmentID)
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, mergedPPMShipment.ExpectedDepartureDate)
		suite.Equal(oldPPMShipment.PickupPostalCode, mergedPPMShipment.PickupPostalCode)
		suite.Equal(oldPPMShipment.DestinationPostalCode, mergedPPMShipment.DestinationPostalCode)
		suite.Equal(oldPPMShipment.SitExpected, mergedPPMShipment.SitExpected)
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

	suite.Run("Passing nil to required fields", func() {
		data := setupShipmentData()
		oldPPMShipment := data.old

		newPPMShipment := models.PPMShipment{
			ExpectedDepartureDate: time.Time{},
			PickupPostalCode:      "",
			DestinationPostalCode: "",
		}

		mergedPPMShipment := mergePPMShipment(newPPMShipment, &oldPPMShipment)

		// Check if old fields aren't updated
		suite.Equal(oldPPMShipment.ExpectedDepartureDate, mergedPPMShipment.ExpectedDepartureDate)
		suite.Equal(oldPPMShipment.PickupPostalCode, mergedPPMShipment.PickupPostalCode)
		suite.Equal(oldPPMShipment.DestinationPostalCode, mergedPPMShipment.DestinationPostalCode)
	})
}
