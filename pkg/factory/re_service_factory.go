package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildReService creates a ReService
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildReService(db *pop.Connection, customs []Customization, traits []Trait) models.ReService {
	customs = setupCustomizations(customs, traits)

	// Find reService assertion and convert to models ReService
	var cReService models.ReService
	if result := findValidCustomization(customs, ReService); result != nil {
		cReService = result.Model.(models.ReService)
		if result.LinkOnly {
			return cReService
		}
	}

	// create reService
	reServiceUUID := uuid.Must(uuid.NewV4())
	reService := models.ReService{
		ID:   reServiceUUID,
		Name: "Test Service",
		Code: models.ReServiceCode("STEST"),
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&reService, cReService)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &reService)
	}

	return reService
}

// FetchOrBuildReServiceByCode tries fetching a ReService using ReServiceCode, then falls back to creating one
func FetchOrBuildReServiceByCode(db *pop.Connection, reServiceCode models.ReServiceCode) models.ReService {
	if db == nil {
		return BuildReService(db, []Customization{
			{
				Model: models.ReService{
					Code: reServiceCode,
				},
			},
		}, nil)
	}

	var reService models.ReService
	err := db.Where("code=$1", reServiceCode).First(&reService)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return reService
	}

	return BuildReService(db, []Customization{
		{
			Model: models.ReService{
				Code: reServiceCode,
			},
		},
	}, nil)
}

// BuildReServiceByCode builds ReService using ReServiceCode
func BuildReServiceByCode(db *pop.Connection, reServiceCode models.ReServiceCode) models.ReService {
	return BuildReService(db, []Customization{
		{
			Model: models.ReService{
				Code: reServiceCode,
			},
		},
	}, nil)
}

// BuildDDFSITReService creates the three destination SIT service codes: DDFSIT, DDASIT, DDDSIT. Returns DDFSIT only.
func BuildDDFSITReService(db *pop.Connection) models.ReService {
	reService := FetchOrBuildReServiceByCode(db, models.ReServiceCodeDDFSIT)
	FetchOrBuildReServiceByCode(db, models.ReServiceCodeDDASIT)
	FetchOrBuildReServiceByCode(db, models.ReServiceCodeDDDSIT)
	return reService
}

// BuildDOFSITReService creates the three origin SIT service codes: DOFSIT, DOPSIT, DOASIT. Returns DOFSIT only.
func BuildDOFSITReService(db *pop.Connection) models.ReService {
	reService := FetchOrBuildReServiceByCode(db, models.ReServiceCodeDOFSIT)
	FetchOrBuildReServiceByCode(db, models.ReServiceCodeDOASIT)
	FetchOrBuildReServiceByCode(db, models.ReServiceCodeDOPSIT)
	return reService
}
