package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Creates a State for Beverly Hills, California 90210
func BuildState(db *pop.Connection, customs []Customization, traits []Trait) models.State {
	customs = setupCustomizations(customs, traits)

	var cState models.State
	if result := findValidCustomization(customs, UsPostRegionCity); result != nil {
		cState = result.Model.(models.State)
		if result.LinkOnly {
			return cState
		}
	}

	state := models.State{
		ID:        uuid.Must(uuid.NewV4()),
		State:     "CA",
		StateName: "CALIFORNIA",
		IsOconus:  false,
	}

	testdatagen.MergeModels(&state, cState)

	if db != nil {
		mustCreate(db, &state)
	}

	return state
}

// Creates a default State for Beverly Hills, California 90210
func BuildDefaultState(db *pop.Connection) models.State {
	return BuildState(db, nil, nil)
}

// FetchOrBuildState tries fetching a State
func FetchOrBuildState(db *pop.Connection, customs []Customization, traits []Trait) models.State {
	if db == nil {
		return BuildState(db, customs, traits)
	}

	customs = setupCustomizations(customs, traits)

	var cState models.State
	if result := findValidCustomization(customs, State); result != nil {
		cState = result.Model.(models.State)
		if result.LinkOnly {
			return cState
		}
	}

	var state models.State
	if !cState.ID.IsNil() {
		err := db.Where("ID = $1", cState.ID).First(&state)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return state
		}
	}

	return BuildState(db, customs, traits)
}
