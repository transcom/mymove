package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFindDutyLocations() {
	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.MustSave(&address)

	location1 := models.DutyLocation{
		Name:      "Fort Bragg",
		AddressID: address.ID,
	}
	suite.MustSave(&location1)

	location2 := models.DutyLocation{
		Name:      "Fort Belvoir",
		AddressID: address.ID,
	}
	suite.MustSave(&location2)

	location3 := models.DutyLocation{
		Name:      "Davis Monthan AFB",
		AddressID: address.ID,
	}
	suite.MustSave(&location3)

	location4 := models.DutyLocation{
		Name:      "JB Elmendorf-Richardson",
		AddressID: address.ID,
	}
	suite.MustSave(&location4)

	location5 := models.DutyLocation{
		Name:      "NAS Fallon",
		AddressID: address.ID,
	}
	suite.MustSave(&location5)

	s5 := models.DutyLocationName{
		Name:           "Naval Air Station Fallon",
		DutyLocationID: location5.ID,
	}
	suite.MustSave(&s5)

	location6 := models.DutyLocation{
		Name:      "NAS Fort Worth JRB",
		AddressID: address.ID,
	}
	suite.MustSave(&location6)
	s6 := models.DutyLocationName{
		Name:           "Naval Air Station Fort Worth Joint Reserve Base",
		DutyLocationID: location6.ID,
	}
	suite.MustSave(&s6)

	address2 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "23456",
	}
	suite.MustSave(&address2)

	location7 := models.DutyLocation{
		Name:      "Very Long City Name, OH 23456",
		AddressID: address2.ID,
	}
	suite.MustSave(&location7)

	tests := []struct {
		query         string
		dutyLocations []string
	}{
		{query: "fort", dutyLocations: []string{"Fort Bragg", "Fort Belvoir", "NAS Fort Worth JRB", "NAS Fallon"}},
		{query: "ft", dutyLocations: []string{"Fort Bragg", "NAS Fallon", "Fort Belvoir", "NAS Fort Worth JRB"}},
		{query: "ft be", dutyLocations: []string{"Fort Belvoir", "Fort Bragg", "NAS Fallon", "NAS Fort Worth JRB"}},
		{query: "davis-mon", dutyLocations: []string{"Davis Monthan AFB", "NAS Fallon", "JB Elmendorf-Richardson"}},
		{query: "jber", dutyLocations: []string{"JB Elmendorf-Richardson", "NAS Fort Worth JRB"}},
		{query: "naval air", dutyLocations: []string{"NAS Fallon", "NAS Fort Worth JRB", "Very Long City Name, OH 23456", "Fort Belvoir", "Davis Monthan AFB"}},
		{query: "zzzzz", dutyLocations: []string{}},
		{query: "23456", dutyLocations: []string{"Very Long City Name, OH 23456"}},
	}

	for _, ts := range tests {
		dutyLocations, err := models.FindDutyLocations(suite.DB(), ts.query)
		suite.NoError(err)
		suite.Equal(len(dutyLocations), len(ts.dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
		for i, dutyLocation := range dutyLocations {
			suite.Equal(dutyLocation.Name, ts.dutyLocations[i], "Duty locations don't match order: %s", ts.query)
		}
	}
}

func (suite *ModelSuite) Test_DutyLocationValidations() {
	location := &models.DutyLocation{}

	var expErrors = map[string][]string{
		"name":       {"Name can not be blank."},
		"address_id": {"AddressID can not be blank."},
	}

	suite.verifyValidationErrors(location, expErrors)
}
func (suite *ModelSuite) Test_FetchDutyLocationTransportationOffice() {
	t := suite.T()
	dutyLocation := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())

	office, err := models.FetchDutyLocationTransportationOffice(suite.DB(), dutyLocation.ID)
	if err != nil {
		t.Errorf("Find transportation office error: %v", err)
	}

	if office.PhoneLines[0].Number != "(510) 555-5555" {
		t.Error("phone number should be loaded")
	}

}
