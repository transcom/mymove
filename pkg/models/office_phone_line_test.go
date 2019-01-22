package models_test

import . "github.com/transcom/mymove/pkg/models"

func (suite *ModelSuite) Test_OfficePhoneLineInstantiation() {
	phoneLine := &OfficePhoneLine{}
	expErrors := map[string][]string{
		"number":                   {"Number can not be blank."},
		"type":                     {"Type is not in the list [voice, fax]."},
		"transportation_office_id": {"TransportationOfficeID can not be blank."},
	}
	suite.verifyValidationErrors(phoneLine, expErrors)
}

func (suite *ModelSuite) Test_BasicOfficePhoneLine() {
	office := CreateTestShippingOffice(suite)
	infoLine := OfficePhoneLine{
		TransportationOfficeID: office.ID,
		Number:                 "(907) 555-1212",
		Label:                  StringPointer("Information Only"),
		Type:                   "voice",
	}

	suite.MustSave(&infoLine)
	suite.False(infoLine.IsDsnNumber)

	faxLine := OfficePhoneLine{
		TransportationOfficeID: office.ID,
		Number:                 "555 12345",
		Label:                  StringPointer("Secure Fax"),
		Type:                   "fax",
		IsDsnNumber:            true,
	}

	suite.MustSave(&faxLine)
	suite.True(faxLine.IsDsnNumber)
	var loadedOffice TransportationOffice
	err := suite.DB().Eager().Find(&loadedOffice, office.ID)
	suite.Nil(err, "loading office with phone lines")
	suite.Equal(2, len(loadedOffice.PhoneLines))
}
