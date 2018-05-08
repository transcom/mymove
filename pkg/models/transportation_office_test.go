package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_TransportationOfficeInstantiation() {
	office := &TransportationOffice{}
	expErrors := map[string][]string{
		"name": {"Name can not be blank."},
		"address": {"Address.StreetAddress1 can not be blank.",
			"Address.City can not be blank.",
			"Address.State can not be blank.",
			"Address.PostalCode can not be blank."},
	}
	suite.verifyValidationErrors(office, expErrors)
}

func NewTestShippingOffice() TransportationOffice {
	return TransportationOffice{
		Name: "JPSO Supreme",
		Address: Address{
			StreetAddress1: "123 washington Ave",
			City:           "Springfield",
			State:          "AK",
			PostalCode:     "99515"},
		Latitude:  61.1262383,
		Longitude: -149.9212882,
		Hours:     StringPointer("0900-1800 Mon-Sat"),
	}
}

func (suite *ModelSuite) Test_BasicShippingOffice() {
	office := NewTestShippingOffice()
	verrs, err := office.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(office.ID, "basic office ID")
	suite.NotNil(office.Address.ID, "didn't save address")
}

func (suite *ModelSuite) Test_TransportationOffice() {
	jppso := NewTestShippingOffice()
	ppo := TransportationOffice{
		Name:           "Best PPO of the North",
		ShippingOffice: &jppso,
		Address: Address{
			StreetAddress1: "456 Lincoln St",
			City:           "Sitka",
			State:          "AK",
			PostalCode:     "99835"},
		Latitude:  57.0512403,
		Longitude: -135.332707,
		Services:  StringPointer("Moose Shipping, Personal Goods"),
	}
	verrs, err := ppo.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(ppo.ID)
	suite.NotNil(ppo.ShippingOffice)
	suite.NotNil(ppo.ShippingOffice.ID)
	suite.NotNil(ppo.ShippingOffice.Address.ID)
}
