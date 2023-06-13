package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestFindDutyLocations() {
	location1 := models.DutyLocation{
		Name:       "Fort Bragg",
		City:       "city",
		State:      "state",
		PostalCode: "12345",
		Country:    "United States",
	}
	suite.MustSave(&location1)

	location2 := models.DutyLocation{
		Name:       "Fort Belvoir",
		City:       "city",
		State:      "state",
		PostalCode: "12345",
		Country:    "United States",
	}
	suite.MustSave(&location2)

	location3 := models.DutyLocation{
		Name:       "Davis Monthan AFB",
		City:       "city",
		State:      "state",
		PostalCode: "12345",
		Country:    "United States",
	}
	suite.MustSave(&location3)

	location4 := models.DutyLocation{
		Name:       "JB Elmendorf-Richardson",
		City:       "city",
		State:      "state",
		PostalCode: "12345",
		Country:    "United States",
	}
	suite.MustSave(&location4)

	location5 := models.DutyLocation{
		Name:       "NAS Fallon",
		City:       "city",
		State:      "state",
		PostalCode: "12345",
		Country:    "United States",
	}
	suite.MustSave(&location5)

	s5 := models.DutyLocationName{
		Name:           "Naval Air Station Fallon",
		DutyLocationID: location5.ID,
	}
	suite.MustSave(&s5)

	location6 := models.DutyLocation{
		Name:       "NAS Fort Worth JRB",
		City:       "city",
		State:      "state",
		PostalCode: "12345",
		Country:    "United States",
	}
	suite.MustSave(&location6)
	s6 := models.DutyLocationName{
		Name:           "Naval Air Station Fort Worth Joint Reserve Base",
		DutyLocationID: location6.ID,
	}
	suite.MustSave(&s6)

	location7 := models.DutyLocation{
		Name:       "Very Long City Name, OH 23456",
		City:       "city",
		State:      "state",
		PostalCode: "23456",
		Country:    "United States",
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
		suite.NoError(err, ts.query)
		suite.Require().Equal(len(dutyLocations), len(ts.dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
		for i, dutyLocation := range dutyLocations {
			suite.Equal(dutyLocation.Name, ts.dutyLocations[i], "Duty locations don't match order: %s", ts.query)
		}
	}
}

func (suite *ModelSuite) Test_DutyLocationValidations() {
	location := &models.DutyLocation{}

	var expErrors = map[string][]string{
		"name":        {"Name can not be blank."},
		"city":        {"City can not be blank."},
		"state":       {"State can not be blank."},
		"postal_code": {"PostalCode can not be blank."},
		"country":     {"Country can not be blank."},
	}

	suite.verifyValidationErrors(location, expErrors)
}
func (suite *ModelSuite) Test_FetchDutyLocationTransportationOffice() {
	t := suite.T()

	suite.Run("fetches duty location with transportation office", func() {
		dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())

		office, err := models.FetchDutyLocationTransportationOffice(suite.DB(), dutyLocation.ID)
		if err != nil {
			t.Errorf("Find transportation office error: %v", err)
		}

		if office.PhoneLines[0].Number != "(510) 555-5555" {
			t.Error("phone number should be loaded")
		}
	})

	suite.Run("if duty location does not have a transportation office, it throws ErrFetchNotFound error and returns and empty office", func() {
		dutyLocationWithoutTransportationOffice := factory.BuildDutyLocationWithoutTransportationOffice(suite.DB(), nil, nil)

		suite.Equal(uuid.Nil, dutyLocationWithoutTransportationOffice.TransportationOffice.ID)

		office, err := models.FetchDutyLocationTransportationOffice(suite.DB(), dutyLocationWithoutTransportationOffice.ID)
		suite.Error(err)
		suite.IsType(models.ErrFetchNotFound, err)
		suite.Equal(models.TransportationOffice{}, office)
	})
}

func (suite *ModelSuite) Test_FetchDutyLocationWithTransportationOffice() {
	suite.Run("fetches duty location with transportation office", func() {
		dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
		officePhoneLine := dutyLocation.TransportationOffice.PhoneLines[0].Number
		dutyLocationFromDB, err := models.FetchDutyLocationWithTransportationOffice(suite.DB(), dutyLocation.ID)
		suite.NoError(err)
		suite.Equal(dutyLocation.TransportationOfficeID, dutyLocationFromDB.TransportationOfficeID)
		suite.Equal(officePhoneLine, dutyLocationFromDB.TransportationOffice.PhoneLines[0].Number)
	})

	suite.Run("if duty location does not have a transportation office, it will still return the duty location without throwing an error", func() {
		dutyLocation := factory.BuildDutyLocationWithoutTransportationOffice(suite.DB(), nil, nil)
		dutyLocationFromDB, err := models.FetchDutyLocationWithTransportationOffice(suite.DB(), dutyLocation.ID)
		suite.NoError(err)
		suite.Nil(dutyLocationFromDB.TransportationOfficeID)
	})
}
