package models_test

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"time"
)

func (suite *ModelSuite) TestBasicFuelEIADieselPriceInstantiation() {
	now := time.Now()
	id := uuid.Must(uuid.NewV4())
	newFuelPrice := models.FuelEIADieselPrice{
		ID:                          id,
		PubDate:                     now,
		RateStartDate:               now,
		RateEndDate:                 now.AddDate(0, 1, -1),
		EIAPricePerGallonMillicents: 320700,
		BaselineRate:                6,
	}

	verrs, err := suite.db.ValidateAndCreate(&newFuelPrice)

	suite.Nil(err, "Error writing to the db.")
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestEmptyFuelEIADieselPriceInstantiation() {
	newFuelPrice := models.FuelEIADieselPrice{}

	expErrors := map[string][]string{
		"pub_date":                          {"PubDate can not be blank."},
		"rate_start_date":                   {"RateStartDate can not be blank."},
		"rate_end_date":                     {"RateEndDate can not be blank."},
		"e_i_a_price_per_gallon_millicents": {"0 is not greater than 0."},
		"baseline_rate":                     {"0 is not greater than 0."},
	}
	suite.verifyValidationErrors(&newFuelPrice, expErrors)
}

func (suite *ModelSuite) TestMakeFuelEIADieselPrices() {
	testdatagen.MakeFuelEIADieselPrices(suite.db)
}

func (suite *ModelSuite) TestMakeFuelEIADieselPriceForDate() {
	rateStartDate := time.Date(2017, time.July, 15, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.RateStartDate = rateStartDate
	assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents = unit.Millicents(695700)
	shipmentDate := assertions.FuelEIADieselPrice.RateStartDate.AddDate(0, 0, 10)

	testdatagen.MakeFuelEIADieselPriceForDate(suite.db, shipmentDate, assertions)
}

func (suite *ModelSuite) TestMakeDefaultFuelEIADieselPriceForDate() {
	shipmentDate := time.Now()
	testdatagen.MakeDefaultFuelEIADieselPriceForDate(suite.db, shipmentDate)
}
