package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildSITDurationUpdate creates an SITDurationUpdate
// It builds
//   - MTOShipment and associated set relationships
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildSITDurationUpdate(db *pop.Connection, customs []Customization, traits []Trait) models.SITDurationUpdate {
	customs = setupCustomizations(customs, traits)

	// Find SITDurationUpdate Customization and extract the custom SITDurationUpdate
	var cSITDurationUpdate models.SITDurationUpdate
	if result := findValidCustomization(customs, SITDurationUpdate); result != nil {
		cSITDurationUpdate = result.Model.(models.SITDurationUpdate)
		if result.LinkOnly {
			return cSITDurationUpdate
		}
	}

	shipment := BuildMTOShipment(db, customs, traits)

	// Create default SITDurationUpdate
	SITDurationUpdate := models.SITDurationUpdate{
		MTOShipment:   shipment,
		MTOShipmentID: shipment.ID,
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		RequestedDays: 45,
		Status:        models.SITExtensionStatusPending,
	}

	// Overwrite default values with those from custom SITDurationUpdate
	testdatagen.MergeModels(&SITDurationUpdate, cSITDurationUpdate)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &SITDurationUpdate)
	}

	return SITDurationUpdate
}

// ------------------------
//      TRAITS
// ------------------------

func GetTraitApprovedSITDurationUpdate() []Customization {
	return []Customization{
		{
			Model: models.SITDurationUpdate{
				RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
				Status:        models.SITExtensionStatusApproved,
				ApprovedDays:  models.IntPointer(100),
				DecisionDate:  models.TimePointer(time.Now()),
				RequestedDays: 90,
			},
		},
	}
}
