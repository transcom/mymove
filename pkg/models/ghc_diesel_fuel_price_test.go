package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGHCDieselFuelPriceInstantiation() {
	ghcDieselFuelPrice := &GHCDieselFuelPrice{}

	expectedErrors := map[string][]string{
		"fuel_price_in_millicents": {"FuelPriceInMillicents can not be blank."},
		"publication_date":         {"PublicationDate can not be blank."},
	}

	suite.verifyValidationErrors(ghcDieselFuelPrice, expectedErrors)
}

func (suite *ModelSuite) TestGHCDieselFuelPriceUniqueness() {
	t := suite.T()
	ghcDieselFuelPrice := &GHCDieselFuelPrice{
		FuelPriceInMillicents: 500000,
		PublicationDate:       time.Now(),
	}

	if verrs, err := suite.DB().ValidateAndCreate(ghcDieselFuelPrice); err != nil || verrs.HasAny() {
		t.Errorf("Didn't create GHC Diesel Fuel Price: %s", err)
	}

	anotherGHCDieselFuelPrice := &GHCDieselFuelPrice{
		FuelPriceInMillicents: 100,
		PublicationDate:       time.Now(),
	}

	_, err := suite.DB().ValidateAndCreate(anotherGHCDieselFuelPrice)

	suite.Error(err)
}
