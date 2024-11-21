package models_test

import (
	"fmt"
	"slices"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
)

func (suite *ModelSuite) TestFindDutyLocations() {
	addressCreator := address.NewAddressCreator()
	newAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	createdAddress, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newAddress)
	suite.NoError(err)

	location1 := models.DutyLocation{
		Name:      "Fort Bragg",
		AddressID: createdAddress.ID,
	}
	suite.MustSave(&location1)

	location2 := models.DutyLocation{
		Name:      "Fort Belvoir",
		AddressID: createdAddress.ID,
	}
	suite.MustSave(&location2)

	location3 := models.DutyLocation{
		Name:      "Davis Monthan AFB",
		AddressID: createdAddress.ID,
	}
	suite.MustSave(&location3)

	location4 := models.DutyLocation{
		Name:      "JB Elmendorf-Richardson",
		AddressID: createdAddress.ID,
	}
	suite.MustSave(&location4)

	location5 := models.DutyLocation{
		Name:      "NAS Fallon",
		AddressID: createdAddress.ID,
	}
	suite.MustSave(&location5)

	s5 := models.DutyLocationName{
		Name:           "Naval Air Station Fallon",
		DutyLocationID: location5.ID,
	}
	suite.MustSave(&s5)

	location6 := models.DutyLocation{
		Name:      "NAS Fort Worth JRB",
		AddressID: createdAddress.ID,
	}
	suite.MustSave(&location6)
	s6 := models.DutyLocationName{
		Name:           "Naval Air Station Fort Worth Joint Reserve Base",
		DutyLocationID: location6.ID,
	}
	suite.MustSave(&s6)

	newAddress2 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "23456",
	}
	createdAddress2, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newAddress2)
	suite.NoError(err)

	location7 := models.DutyLocation{
		Name:      "Very Long City Name, OH 23456",
		AddressID: createdAddress2.ID,
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
		suite.Require().Equal(len(dutyLocations), len(ts.dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
		for i, dutyLocation := range dutyLocations {
			suite.Equal(dutyLocation.Name, ts.dutyLocations[i], "Duty locations don't match order: %s", ts.query)
		}
	}
}

func (suite *ModelSuite) TestFindDutyLocationExcludeStates() {
	addressCreator := address.NewAddressCreator()
	newAKAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "AK",
		PostalCode:     "12345",
	}
	createdAddress1, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newAKAddress)
	suite.NoError(err)

	newHIAddress := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "HI",
		PostalCode:     "12345",
	}
	createdAddress2, err := addressCreator.CreateAddress(suite.AppContextForTest(), &newHIAddress)
	suite.NoError(err)

	location1 := models.DutyLocation{
		Name:      "Fort Test 1",
		AddressID: createdAddress1.ID,
	}
	suite.MustSave(&location1)

	location2 := models.DutyLocation{
		Name:      "Fort Test 2",
		AddressID: createdAddress2.ID,
	}
	suite.MustSave(&location2)

	tests := []struct {
		query         string
		dutyLocations []string
	}{
		{query: "fort test", dutyLocations: []string{"Fort Test 1", "Fort Test 2"}},
	}

	statesToExclude := make([]string, 0)
	statesToExclude = append(statesToExclude, "AK")
	statesToExclude = append(statesToExclude, "HI")

	for _, ts := range tests {
		dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, statesToExclude)
		suite.NoError(err)
		suite.Require().Equal(0, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
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

func (suite *ModelSuite) Test_SearchDutyLocations_Exclude_Not_Active_Oconus() {
	createContract := func(appCtx appcontext.AppContext, contractCode string, contractName string) (*models.ReContract, error) {
		// See if contract code already exists.
		exists, err := appCtx.DB().Where("code = ?", contractCode).Exists(&models.ReContract{})
		if err != nil {
			return nil, fmt.Errorf("could not determine if contract code [%s] existed: %w", contractCode, err)
		}
		if exists {
			return nil, fmt.Errorf("the provided contract code [%s] already exists", contractCode)
		}

		// Contract code is new; insert it.
		contract := models.ReContract{
			Code: contractCode,
			Name: contractName,
		}
		verrs, err := appCtx.DB().ValidateAndSave(&contract)
		if verrs.HasAny() {
			return nil, fmt.Errorf("validation errors when saving contract [%+v]: %w", contract, verrs)
		}
		if err != nil {
			return nil, fmt.Errorf("could not save contract [%+v]: %w", contract, err)
		}
		return &contract, nil
	}

	setupDataForOconusSearchCounselingOffice := func(contract models.ReContract, postalCode string, gbloc string, dutyLocationName string, transportationName string, isOconusRateAreaActive bool) (models.ReRateArea, models.OconusRateArea, models.UsPostRegionCity, models.DutyLocation) {
		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := models.ReRateArea{
			ID:         uuid.Must(uuid.NewV4()),
			ContractID: contract.ID,
			IsOconus:   true,
			Code:       rateAreaCode,
			Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
			Contract:   contract,
		}
		verrs, err := suite.DB().ValidateAndCreate(&rateArea)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		us_country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), postalCode)
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		oconusRateArea := models.OconusRateArea{
			ID:                 uuid.Must(uuid.NewV4()),
			RateAreaId:         rateArea.ID,
			CountryId:          us_country.ID,
			UsPostRegionCityId: usprc.ID,
			Active:             isOconusRateAreaActive,
		}
		verrs, err = suite.DB().ValidateAndCreate(&oconusRateArea)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		address := models.Address{
			StreetAddress1:     "n/a",
			City:               "SomeCity",
			State:              "AK",
			PostalCode:         postalCode,
			County:             "SomeCounty",
			IsOconus:           models.BoolPointer(true),
			UsPostRegionCityId: &usprc.ID,
			CountryId:          models.UUIDPointer(us_country.ID),
		}
		suite.MustSave(&address)

		origDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name:                       dutyLocationName,
					AddressID:                  address.ID,
					ProvidesServicesCounseling: true,
				},
			},
			{
				Model: models.TransportationOffice{
					Name:             transportationName,
					Gbloc:            gbloc,
					ProvidesCloseout: true,
				},
			},
		}, nil)
		suite.MustSave(&origDutyLocation)

		found_duty_location, _ := models.FetchDutyLocation(suite.DB(), origDutyLocation.ID)

		return rateArea, oconusRateArea, *usprc, found_duty_location
	}

	const fairbanksAlaskaPostalCode = "99790"
	const anchorageAlaskaPostalCode = "99502"
	testContractName := "Test_search_duty_location"
	testContractCode := "Test_search_duty_location_Code"
	testGbloc := "ABCD"
	testTransportationName := "TEST - PPO"
	testDutyLocationName := "TEST Duty Location"
	testTransportationName2 := "TEST - PPO 2"
	testDutyLocationName2 := "TEST Duty Location 2"

	suite.Run("one active onconus rateArea duty location and one not active oconus rate area duty location should return 1", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// active duty location
		_, oconusRateArea, _, dutyLocation := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, true)

		// not active duty location
		_, oconusRateArea2, _, _ := setupDataForOconusSearchCounselingOffice(*contract, anchorageAlaskaPostalCode, testGbloc, testDutyLocationName2, testTransportationName2, false)

		suite.True(oconusRateArea.Active)
		suite.False(oconusRateArea2.Active)

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations test", dutyLocations: []string{testDutyLocationName}},
		}

		expectedDutyLocationNames := []string{dutyLocation.Name}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(1, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
			for _, o := range dutyLocations {
				suite.True(slices.Contains(expectedDutyLocationNames, o.Name))
			}
		}
	})

	suite.Run("two active onconus rateArea duty locations should return 2", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// active duty location
		_, oconusRateArea, _, dutyLocation1 := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, true)

		// active duty location
		_, oconusRateArea2, _, dutyLocation2 := setupDataForOconusSearchCounselingOffice(*contract, anchorageAlaskaPostalCode, testGbloc, testDutyLocationName2, testTransportationName2, true)

		suite.True(oconusRateArea.Active)
		suite.True(oconusRateArea2.Active)

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations test", dutyLocations: []string{testDutyLocationName}},
		}

		expectedDutyLocationNames := []string{dutyLocation1.Name, dutyLocation2.Name}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(2, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
			for _, o := range dutyLocations {
				suite.True(slices.Contains(expectedDutyLocationNames, o.Name))
			}
		}
	})

	suite.Run("two inactive onconus rateArea duty locations should return 0", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// active duty location
		_, oconusRateArea, _, dutyLocation1 := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, false)

		// active duty location
		_, oconusRateArea2, _, dutyLocation2 := setupDataForOconusSearchCounselingOffice(*contract, anchorageAlaskaPostalCode, testGbloc, testDutyLocationName2, testTransportationName2, false)

		suite.False(oconusRateArea.Active)
		suite.False(oconusRateArea2.Active)

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations test", dutyLocations: []string{testDutyLocationName}},
		}

		expectedDutyLocationNames := []string{dutyLocation1.Name, dutyLocation2.Name}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(0, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
			for _, o := range dutyLocations {
				suite.True(slices.Contains(expectedDutyLocationNames, o.Name))
			}
		}
	})

	suite.Run("match on alternative name but exclude", func() {
		alternativeDutyLocationName := "Foobar"

		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// not active duty location
		_, oconusRateArea, _, dutyLocation1 := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, false)

		suite.False(oconusRateArea.Active)

		dutyLocationName := models.DutyLocationName{
			Name:           alternativeDutyLocationName,
			DutyLocationID: dutyLocation1.ID,
			DutyLocation:   dutyLocation1,
		}
		verrs, err := suite.DB().ValidateAndCreate(&dutyLocationName)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations", dutyLocations: []string{alternativeDutyLocationName}}, //search on alt name
		}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(0, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
		}
	})

	suite.Run("match on alternative name when active oconus rateArea", func() {
		alternativeDutyLocationName := "Foobar"

		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// active duty location
		_, oconusRateArea, _, dutyLocation1 := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, true)

		suite.True(oconusRateArea.Active)

		dutyLocationName := models.DutyLocationName{
			Name:           alternativeDutyLocationName,
			DutyLocationID: dutyLocation1.ID,
			DutyLocation:   dutyLocation1,
		}
		verrs, err := suite.DB().ValidateAndCreate(&dutyLocationName)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations", dutyLocations: []string{alternativeDutyLocationName}}, //search on alt name
		}

		expectedDutyLocationNames := []string{dutyLocation1.Name}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(1, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
			for _, o := range dutyLocations {
				suite.True(slices.Contains(expectedDutyLocationNames, o.Name))
			}
		}
	})

	suite.Run("two active onconus rateArea duty locations - search by zip - return match", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// active duty location
		_, oconusRateArea, _, dutyLocation1 := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, true)

		// active duty location
		_, oconusRateArea2, _, dutyLocation2 := setupDataForOconusSearchCounselingOffice(*contract, anchorageAlaskaPostalCode, testGbloc, testDutyLocationName2, testTransportationName2, true)

		suite.True(oconusRateArea.Active)
		suite.True(oconusRateArea2.Active)

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations test", dutyLocations: []string{
				fairbanksAlaskaPostalCode, anchorageAlaskaPostalCode}}, //search by zip
		}

		expectedDutyLocationNames := []string{dutyLocation1.Name, dutyLocation2.Name}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(1, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
			for _, o := range dutyLocations {
				suite.True(slices.Contains(expectedDutyLocationNames, o.Name))
			}
		}
	})

	suite.Run("two non active onconus rateArea duty locations - search by zip - return none", func() {
		contract, err := createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// not active duty location
		_, oconusRateArea, _, _ := setupDataForOconusSearchCounselingOffice(*contract, fairbanksAlaskaPostalCode, testGbloc, testDutyLocationName, testTransportationName, false)

		// not active duty location
		_, oconusRateArea2, _, _ := setupDataForOconusSearchCounselingOffice(*contract, anchorageAlaskaPostalCode, testGbloc, testDutyLocationName2, testTransportationName2, false)

		suite.False(oconusRateArea.Active)
		suite.False(oconusRateArea2.Active)

		tests := []struct {
			query         string
			dutyLocations []string
		}{
			{query: "search oconus rate area duty locations test", dutyLocations: []string{
				fairbanksAlaskaPostalCode, anchorageAlaskaPostalCode}}, //search by zip
		}

		for _, ts := range tests {
			dutyLocations, err := models.FindDutyLocationsExcludingStates(suite.DB(), ts.query, []string{})
			suite.NoError(err)
			suite.Require().Equal(0, len(dutyLocations), "Wrong number of duty locations returned from query: %s", ts.query)
		}
	})
}
