package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const defaultServiceCode = models.ReServiceCode("DLH")

func fetchDefaultReService(db *pop.Connection) models.ReService {
	var reService models.ReService
	err := db.Where("code = $1", defaultServiceCode).First(&reService)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}
	return reService
}

// FetchReService tries fetching a ReService using ReServiceCode, then falls back to creating one
func FetchReService(db *pop.Connection, customs []Customization, traits []Trait) models.ReService {

	customs = setupCustomizations(customs, traits)

	// Find reService assertion and convert to models ReService
	var cReService models.ReService
	if result := findValidCustomization(customs, ReService); result != nil {
		cReService = result.Model.(models.ReService)
		if result.LinkOnly {
			return cReService
		}
	}

	if db == nil {
		defaultReService := models.ReService{
			Name: "Domestic linehaul",
			Code: defaultServiceCode,
		}
		// Overwrite values with those from assertions
		testdatagen.MergeModels(&defaultReService, cReService)
		return defaultReService
	}

	var reService models.ReService
	if !cReService.ID.IsNil() {
		err := db.Where("ID = $1", cReService.ID).First(&reService)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		} else if err == nil {
			return reService
		}
	}

	// search for the default code if one is not provided
	reServiceCode := defaultServiceCode

	if cReService.Code.String() != "" {
		reServiceCode = cReService.Code
	}
	err := db.Where("code = $1", reServiceCode).First(&reService)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return reService
	}

	return fetchDefaultReService(db)
}

func FetchReServiceByCode(db *pop.Connection, reServiceCode models.ReServiceCode) models.ReService {
	return FetchReService(db, []Customization{
		{
			Model: models.ReService{
				Code: reServiceCode,
			},
		},
	}, nil)
}
