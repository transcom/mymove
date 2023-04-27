package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildMTOAgent(db *pop.Connection, customs []Customization, traits []Trait) models.MTOAgent {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cMTOAgent models.MTOAgent
	if result := findValidCustomization(customs, MTOAgent); result != nil {
		cMTOAgent = result.Model.(models.MTOAgent)

		if result.LinkOnly {
			return cMTOAgent
		}
	}

	shipment := BuildMTOShipment(db, customs, traits)

	mtoAgent := models.MTOAgent{
		MTOShipment:   shipment,
		MTOShipmentID: shipment.ID,
		FirstName:     models.StringPointer("Jason"),
		LastName:      models.StringPointer("Ash"),
		Email:         models.StringPointer("jason.ash@example.com"),
		Phone:         models.StringPointer("202-555-9301"),
		MTOAgentType:  models.MTOAgentReleasing,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&mtoAgent, cMTOAgent)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &mtoAgent)
	}

	return mtoAgent
}
