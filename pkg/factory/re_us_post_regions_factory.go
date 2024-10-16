package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Creates a UsPostRegion for Beverly Hills, California 90210
func BuildUsPostRegion(db *pop.Connection, customs []Customization, traits []Trait) models.UsPostRegion {
	customs = setupCustomizations(customs, traits)

	var cUsPostRegion models.UsPostRegion
	if result := findValidCustomization(customs, City); result != nil {
		cUsPostRegion = result.Model.(models.UsPostRegion)
		if result.LinkOnly {
			return cUsPostRegion
		}
	}

	usPostRegion := models.UsPostRegion{
		ID:        uuid.Must(uuid.NewV4()),
		UsprZipID: "90210",
		Zip3:      "902",
	}

	// Find/create the State if customization is provided
	var state models.State
	if result := findValidCustomization(customs, State); result != nil {
		state = BuildState(db, customs, nil)
	} else {
		state = FetchOrBuildState(db, []Customization{
			{
				Model: models.State{
					State:     "FL",
					StateName: "FLORIDA",
					IsOconus:  false,
				},
			},
		}, nil)
	}

	usPostRegion.State = state
	usPostRegion.StateId = state.ID

	testdatagen.MergeModels(&usPostRegion, cUsPostRegion)

	if db != nil {
		mustCreate(db, &usPostRegion)
	}

	return usPostRegion
}

// Creates a default UsPostRegion for Beverly Hills, California 90210
func BuildDefaultUsPostRegion(db *pop.Connection) models.UsPostRegion {
	return BuildUsPostRegion(db, nil, nil)
}
