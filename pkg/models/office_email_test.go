package models_test

import . "github.com/transcom/mymove/pkg/models"

func (suite *ModelSuite) Test_OfficeEmailInstantiation() {
	office := &OfficeEmail{}
	expErrors := map[string][]string{
		"email":                    {"Email can not be blank."},
		"transportation_office_id": {"TransportationOfficeID can not be blank."},
	}
	suite.verifyValidationErrors(office, expErrors)
}
func (suite *ModelSuite) Test_BasicOfficeEmail() {
	office := CreateTestShippingOffice(suite)
	suite.MustSave(&office)

	infoEmail := OfficeEmail{
		TransportationOfficeID: office.ID,
		Email:                  "info@ak_jppso.government.gov",
		Label:                  StringPointer("Information Only"),
	}

	suite.MustSave(&infoEmail)
	suite.NotNil(infoEmail.ID)

	appointmentsEmail := OfficeEmail{
		TransportationOfficeID: office.ID,
		Email:                  "appointments@ak_jppso.government.gov",
	}

	suite.MustSave(&appointmentsEmail)
	suite.NotNil(infoEmail.ID)

	var eagerOffice TransportationOffice
	err := suite.DB().Eager().Find(&eagerOffice, office.ID)
	suite.Nil(err, "Loading office with emails")
	suite.Equal(2, len(eagerOffice.Emails), "Total email count")
}
