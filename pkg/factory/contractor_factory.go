package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildContractor creates a single Contractor
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildContractor(db *pop.Connection, customs []Customization, traits []Trait) models.Contractor {
	customs = setupCustomizations(customs, traits)

	var cContractor models.Contractor
	if result := findValidCustomization(customs, Contractor); result != nil {
		cContractor = result.Model.(models.Contractor)
		if result.LinkOnly {
			return cContractor
		}
	}

	contractor := models.Contractor{
		Name:           DefaultContractName,
		ContractNumber: DefaultContractNumber,
		Type:           DefaultContractType,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&contractor, cContractor)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &contractor)
	}
	return contractor
}

// FetchOrBuildDefaultContractor tries fetching an existing contractor, then falls back to creating one
func FetchOrBuildDefaultContractor(db *pop.Connection, customs []Customization, traits []Trait) models.Contractor {
	if db == nil {
		return BuildContractor(db, customs, traits)
	}

	var contractor models.Contractor
	err := db.Q().Where(`contract_number=$1`, DefaultContractNumber).First(&contractor)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return contractor
	}

	return BuildContractor(db, customs, traits)

}
