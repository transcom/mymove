package models_test

import (
	"github.com/gobuffalo/uuid"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_OfficeUserInstantiation() {
	user := &OfficeUser{}
	expErrors := map[string][]string{
		"given_name":  {"GivenName can not be blank."},
		"family_name": {"FamilyName can not be blank."},
		"telephone":   {"Telephone can not be blank."},
		"email":       {"Email can not be blank."},
		"user": {"User.LoginGovEmail can not be blank.",
			"User.LoginGovUUID can not be blank."},
		"transportation_office": {"TransportationOffice.Name can not be blank.",
			"TransportationOffice.Address.StreetAddress1 can not be blank.",
			"TransportationOffice.Address.City can not be blank.",
			"TransportationOffice.Address.State can not be blank.",
			"TransportationOffice.Address.PostalCode can not be blank."},
	}
	suite.verifyValidationErrors(user, expErrors)
}

func (suite *ModelSuite) Test_BasicOfficeUser() {
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"
	user := OfficeUser{
		User: User{
			LoginGovUUID:  fakeUUID,
			LoginGovEmail: userEmail,
		},
		FamilyName:           "Tester",
		GivenName:            "Sally",
		Email:                "sally.work@government.gov",
		Telephone:            "(907) 555-1212",
		TransportationOffice: NewTestShippingOffice(),
	}
	verrs, err := user.ValidateCreate(suite.db)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(user.ID)
}
