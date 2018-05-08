package models_test

import . "github.com/transcom/mymove/pkg/models"

func (suite *ModelSuite) Test_OfficeEmailInstantiation() {
	office := &OfficeEmail{}
	expErrors := map[string][]string{
		"email": {"Email can not be blank."},
		"transportation_office": {"TransportationOffice.Name can not be blank.",
			"TransportationOffice.Address.StreetAddress1 can not be blank.",
			"TransportationOffice.Address.City can not be blank.",
			"TransportationOffice.Address.State can not be blank.",
			"TransportationOffice.Address.PostalCode can not be blank."},
	}
	suite.verifyValidationErrors(office, expErrors)
}
func (suite *ModelSuite) Test_BasicOfficeEmail() {
	infoEmail := OfficeEmail{
		TransportationOffice: NewTestShippingOffice(),
		Email:                "info@ak_jppso.government.gov",
		Label:                StringPointer("Information Only"),
	}

	verrs, err := infoEmail.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(infoEmail.ID)

	appointmentsEmail := OfficeEmail{
		TransportationOffice: infoEmail.TransportationOffice,
		Email:                "appointments@ak_jppso.government.gov",
	}
	verrs, err = appointmentsEmail.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(infoEmail.ID)
}
