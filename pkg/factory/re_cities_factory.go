package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Creates a City for Beverly Hills, California 90210
func BuildCity(db *pop.Connection, customs []Customization, traits []Trait) models.City {
	customs = setupCustomizations(customs, traits)

	var cCity models.City
	if result := findValidCustomization(customs, City); result != nil {
		cCity = result.Model.(models.City)
		if result.LinkOnly {
			return cCity
		}
	}

	// Check if the state provided already exists in the database
	if db != nil {
		var existingCity models.City
		err := db.Where("city_id = ?", cCity.ID).First(&existingCity)
		if err == nil {
			return existingCity
		}
	}

	city := models.City{
		CityName: "Beverly Hills",
		IsOconus: false,
	}

	// Find/create the Country if customization is provided
	var country models.Country
	if result := findValidCustomization(customs, Country); result != nil {
		country = BuildCountry(db, customs, nil)
	} else {
		country = FetchOrBuildCountry(db, []Customization{
			{
				Model: models.Country{
					Country:     "US",
					CountryName: "UNITED STATES",
				},
			},
		}, nil)
	}

	// Find/create the State if customization is provided
	var state models.State
	if result := findValidCustomization(customs, State); result != nil {
		state = BuildState(db, customs, nil)
	} else {
		state = FetchOrBuildState(db, []Customization{
			{
				Model: models.State{
					State:     "CA",
					StateName: "CALIFORNIA",
					IsOconus:  false,
				},
			},
		}, nil)
	}

	city.Country = country
	city.CountryId = country.ID
	city.State = state
	city.StateId = state.ID

	testdatagen.MergeModels(&city, cCity)

	if db != nil {
		mustCreate(db, &city)
	}

	return city
}

// Creates a default UsPostRegionCity for Beverly Hills, California 90210
func BuildDefaultCity(db *pop.Connection) models.City {
	return BuildCity(db, nil, nil)
}
