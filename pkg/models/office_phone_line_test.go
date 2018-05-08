package models_test

import . "github.com/transcom/mymove/pkg/models"

func (suite *ModelSuite) Test_OfficePhoneLineInstantiation() {
	user := &OfficePhoneLine{}
	expErrors := map[string][]string{
		"number": {"Number can not be blank."},
		"type":   {"Type is not in the list [voice, fax]."},
		"transportation_office": {"TransportationOffice.Name can not be blank.",
			"TransportationOffice.Address.StreetAddress1 can not be blank.",
			"TransportationOffice.Address.City can not be blank.",
			"TransportationOffice.Address.State can not be blank.",
			"TransportationOffice.Address.PostalCode can not be blank."},
	}
	suite.verifyValidationErrors(user, expErrors)
}

func (suite *ModelSuite) Test_BasicOfficePhoneline() {
	infoLine := OfficePhoneLine{
		TransportationOffice: NewTestShippingOffice(),
		Number:               "(907) 555-1212",
		Label:                StringPointer("Information Only"),
		Type:                 "voice",
	}

	verrs, err := infoLine.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(infoLine.ID)
	suite.False(infoLine.IsDsnNumber)

	faxLine := OfficePhoneLine{
		TransportationOffice: infoLine.TransportationOffice,
		Number:               "555 12345",
		Label:                StringPointer("Secure Fax"),
		IsDsnNumber:          true,
	}
	verrs, err = faxLine.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(faxLine.ID)
	suite.True(faxLine.IsDsnNumber)
}
