package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFindDutyStations() {
	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.MustSave(&address)

	station1 := models.DutyStation{
		Name:        "Fort Bragg",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station1)

	station2 := models.DutyStation{
		Name:        "Fort Belvoir",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station2)

	station3 := models.DutyStation{
		Name:        "Davis Monthan AFB",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station3)

	station4 := models.DutyStation{
		Name:        "JB Elmendorf-Richardson",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station4)

	station5 := models.DutyStation{
		Name:        "NAS Fallon",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station5)

	s5 := models.DutyStationName{
		Name:          "Naval Air Station Fallon",
		DutyStationID: station5.ID,
	}
	suite.MustSave(&s5)

	station6 := models.DutyStation{
		Name:        "NAS Fort Worth JRB",
		Affiliation: internalmessages.AffiliationARMY,
		AddressID:   address.ID,
	}
	suite.MustSave(&station6)
	s6 := models.DutyStationName{
		Name:          "Naval Air Station Fort Worth Joint Reserve Base",
		DutyStationID: station6.ID,
	}
	suite.MustSave(&s6)

	tests := []struct {
		query        string
		dutyStations []string
	}{
		{query: "fort", dutyStations: []string{"Fort Bragg", "Fort Belvoir", "NAS Fort Worth JRB", "NAS Fallon"}},
		{query: "ft", dutyStations: []string{"Fort Bragg", "NAS Fallon", "Fort Belvoir", "NAS Fort Worth JRB"}},
		{query: "ft be", dutyStations: []string{"Fort Belvoir", "Fort Bragg", "NAS Fallon", "NAS Fort Worth JRB"}},
		{query: "davis-mon", dutyStations: []string{"Davis Monthan AFB", "NAS Fallon", "JB Elmendorf-Richardson"}},
		{query: "jber", dutyStations: []string{"JB Elmendorf-Richardson", "NAS Fort Worth JRB"}},
		{query: "naval air", dutyStations: []string{"NAS Fallon", "NAS Fort Worth JRB", "Fort Belvoir", "Davis Monthan AFB"}},
		{query: "zzzzz", dutyStations: []string{}},
	}

	for _, ts := range tests {
		dutyStations, err := models.FindDutyStations(suite.DB(), ts.query)
		suite.NoError(err)
		suite.Equal(len(dutyStations), len(ts.dutyStations), "Wrong number of duty stations returned from query: %s", ts.query)
		for i, dutyStation := range dutyStations {
			suite.Equal(dutyStation.Name, ts.dutyStations[i], "Duty stations don't match order: %s", ts.query)
		}
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
