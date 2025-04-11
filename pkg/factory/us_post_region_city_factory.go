package factory

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func FetchOrBuildUsPostRegionCity(db *pop.Connection, customs []Customization, traits []Trait, address *models.Address) models.UsPostRegionCity {
	customs = setupCustomizations(customs, traits)

	var cUsPostRegionCity models.UsPostRegionCity
	if result := findValidCustomization(customs, UsPostRegionCity); result != nil {
		cUsPostRegionCity = result.Model.(models.UsPostRegionCity)
		if result.LinkOnly {
			return cUsPostRegionCity
		}
	}

	usPostRegionCity := models.UsPostRegionCity{
		UsprZipID:          "90210",
		USPostRegionCityNm: "BEVERLY HILLS",
		State:              "CA",
		CtryGencDgphCd:     "US",
		UsprcCountyNm:      "LOS ANGELES",
		UsPostRegionId:     uuid.FromStringOrNil("5a6c650f-f4a9-428a-ae9d-20a251769dc5"),
		CityId:             uuid.FromStringOrNil("d684959a-f59c-4c05-b7c8-0a16df6718aa"),
	}

	if address.PostalCode != "" && address.City != "" {
		usPostRegionCity.USPostRegionCityNm = address.City
		usPostRegionCity.UsprZipID = address.PostalCode
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&usPostRegionCity, cUsPostRegionCity)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {

		fetchedUsPostRegionCity, err := models.FindByZipCodeAndCity(db, usPostRegionCity.UsprZipID, usPostRegionCity.USPostRegionCityNm)
		if err != nil && err != sql.ErrNoRows {
			if errors.Unwrap(err) != nil {
				unWrappedErr := errors.Unwrap(err)
				if unWrappedErr == sql.ErrNoRows {
					mustCreate(db, &usPostRegionCity)
				} else {
					log.Panic(err)
				}
			} else {
				log.Panic(err)
			}
		} else if err == nil && fetchedUsPostRegionCity != nil {
			return *fetchedUsPostRegionCity
		} else {
			mustCreate(db, &usPostRegionCity)
		}

	}

	if usPostRegionCity.ID == uuid.Nil && db == nil {
		usPostRegionCity.ID = uuid.Must(uuid.NewV4())
	}

	return usPostRegionCity
}

func FetchOrBuildUsPostRegionCityForAddress(db *pop.Connection, customs []Customization, traits []Trait, address *models.Address) models.UsPostRegionCity {
	return FetchOrBuildUsPostRegionCity(db, customs, traits, address)
}

func FetchOrBuildUsPostRegionCityDefault(db *pop.Connection, customs []Customization, traits []Trait) models.UsPostRegionCity {
	return FetchOrBuildUsPostRegionCity(db, customs, traits, nil)
}
