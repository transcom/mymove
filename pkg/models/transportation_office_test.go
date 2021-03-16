//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_TransportationOfficeInstantiation() {
	office := &TransportationOffice{}
	expErrors := map[string][]string{
		"name":       {"Name can not be blank."},
		"address_id": {"AddressID can not be blank."},
	}
	suite.verifyValidationErrors(office, expErrors)
}

func CreateTestShippingOffice(suite *ModelSuite) TransportationOffice {
	address := Address{
		StreetAddress1: "123 washington Ave",
		City:           "Springfield",
		State:          "AK",
		PostalCode:     "99515"}
	suite.MustSave(&address)
	office := TransportationOffice{
		Name:      "JPSO Supreme",
		AddressID: address.ID,
		Gbloc:     "BMAF",
		Latitude:  61.1262383,
		Longitude: -149.9212882,
		Hours:     StringPointer("0900-1800 Mon-Sat"),
	}
	suite.MustSave(&office)
	return office
}

func (suite *ModelSuite) Test_BasicShippingOffice() {
	office := CreateTestShippingOffice(suite)
	var loadedOffice TransportationOffice
	suite.DB().Eager().Find(&loadedOffice, office.ID)
	suite.Equal(office.ID, loadedOffice.ID)
	suite.Equal(office.AddressID, loadedOffice.Address.ID)
}

func (suite *ModelSuite) Test_TransportationOffice() {
	jppso := CreateTestShippingOffice(suite)
	ppoAddress := Address{
		StreetAddress1: "456 Lincoln St",
		City:           "Sitka",
		State:          "AK",
		PostalCode:     "99835"}
	suite.MustSave(&ppoAddress)
	ppo := TransportationOffice{
		Name:             "Best PPO of the North",
		ShippingOfficeID: &jppso.ID,
		AddressID:        ppoAddress.ID,
		Gbloc:            "ACQR",
		Latitude:         57.0512403,
		Longitude:        -135.332707,
		Services:         StringPointer("Moose Shipping, Personal Goods"),
	}
	suite.MustSave(&ppo)
	var loadedOffice TransportationOffice
	suite.DB().Eager().Find(&loadedOffice, ppo.ID)
	suite.Equal(ppo.ID, loadedOffice.ID)
	suite.Equal(jppso.ID, loadedOffice.ShippingOffice.ID)
}
