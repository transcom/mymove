package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildState creates a single State.
// Also creates, if not provided:
// - State
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildState(db *pop.Connection, customs []Customization, traits []Trait) models.State {
	customs = setupCustomizations(customs, traits)

	var cState models.State
	if result := findValidCustomization(customs, State); result != nil {
		cState = result.Model.(models.State)
		if result.LinkOnly {
			return cState
		}
	}

	// Check if the state provided already exists in the database
	if db != nil {
		var existingState models.State
		err := db.Where("state = ?", cState.State).First(&existingState)
		if err == nil {
			return existingState
		}
	}

	state := models.State{
		State:     "CA",
		StateName: "California",
		IsOconus:  false,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&state, cState)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &state)
	}
	return state
}

// FetchOrBuildState tries fetching a State using a provided customization, then falls back to creating a default "CA" state
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
	if !state.ID.IsNil() {
		err := db.Where("ID = $1", state.ID).First(&state)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return state
		}
	}

	// Search for the default state code if one is not provided
	defaultStateCode := "CA"
	err := db.Where("state = $1", defaultStateCode).First(&state)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return state
	}

	return BuildState(db, customs, traits)
}
