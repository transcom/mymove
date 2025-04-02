package models_test

import m "github.com/transcom/mymove/pkg/models"

func (suite *ModelSuite) Test_OfficePhoneLineInstantiation() {
	phoneLine := &m.OfficePhoneLine{}
	expErrors := map[string][]string{
		"number":                   {"Number can not be blank."},
		"type":                     {"Type is not in the list [voice, fax]."},
		"transportation_office_id": {"TransportationOfficeID can not be blank."},
	}
	suite.verifyValidationErrors(phoneLine, expErrors, nil)
}

func (suite *ModelSuite) Test_BasicOfficePhoneLine() {
	office := CreateTestShippingOffice(suite)
	infoLine := m.OfficePhoneLine{
		TransportationOfficeID: office.ID,
		Number:                 "(907) 555-1212",
		Label:                  m.StringPointer("Information Only"),
		Type:                   "voice",
	}

	suite.MustSave(&infoLine)
	suite.False(infoLine.IsDsnNumber)

	faxLine := m.OfficePhoneLine{
		TransportationOfficeID: office.ID,
		Number:                 "555 12345",
		Label:                  m.StringPointer("Secure Fax"),
		Type:                   "fax",
		IsDsnNumber:            true,
	}

	suite.MustSave(&faxLine)
	suite.True(faxLine.IsDsnNumber)
	var loadedOffice m.TransportationOffice
	err := suite.DB().Eager().Find(&loadedOffice, office.ID)
	suite.Nil(err, "loading office with phone lines")
	suite.Equal(2, len(loadedOffice.PhoneLines))
}
