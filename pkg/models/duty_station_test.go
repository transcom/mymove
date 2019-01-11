package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFindDutyStations() {
	t := suite.T()

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.MustSave(&address)

	station1 := models.DutyStation{
		Name:        "First Station",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station1)

	station2 := models.DutyStation{
		Name:        "Second Station",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station2)

	stations, err := models.FindDutyStations(suite.DB(), "first")
	if err != nil {
		t.Errorf("Find duty stations error: %v", err)
	}

	if len(stations) != 1 {
		t.Errorf("Should have only got 1 response, got %v", len(stations))
	}

	if stations[0].Name != "First Station" {
		t.Errorf("Station name should have been \"First Station \", got %v", stations[0].Name)
	}

	if stations[0].Address.City != "city" {
		t.Error("Address should have been loaded")
	}
}

func (suite *ModelSuite) Test_DutyStationValidations() {
	station := &models.DutyStation{}

	var expErrors = map[string][]string{
		"name":        {"Name can not be blank."},
		"affiliation": {"Affiliation can not be blank."},
		"address_id":  {"AddressID can not be blank."},
	}

	suite.verifyValidationErrors(station, expErrors)
}
func (suite *ModelSuite) Test_FetchDutyStationTransportationOffice() {
	t := suite.T()
	dutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())

	office, err := models.FetchDutyStationTransportationOffice(suite.DB(), dutyStation.ID)
	if err != nil {
		t.Errorf("Find transportation office error: %v", err)
	}

	if office.PhoneLines[0].Number != "(510) 555-5555" {
		t.Error("phone number should be loaded")
	}

}
