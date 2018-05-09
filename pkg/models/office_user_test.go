package models_test

import (
	"github.com/gobuffalo/uuid"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_OfficeUserInstantiation() {
	user := &OfficeUser{}
	expErrors := map[string][]string{
		"given_name":               {"GivenName can not be blank."},
		"family_name":              {"FamilyName can not be blank."},
		"telephone":                {"Telephone can not be blank."},
		"email":                    {"Email can not be blank."},
		"user_id":                  {"UserID can not be blank."},
		"transportation_office_id": {"TransportationOfficeID can not be blank."},
	}
	suite.verifyValidationErrors(user, expErrors)
}

func (suite *ModelSuite) Test_BasicOfficeUser() {
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"
	sally := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}
	suite.mustSave(&sally)
	office := CreateTestShippingOffice(suite)

	user := OfficeUser{
		FamilyName:             "Tester",
		GivenName:              "Sally",
		Email:                  "sally.work@government.gov",
		Telephone:              "(907) 555-1212",
		UserID:                 sally.ID,
		TransportationOfficeID: office.ID,
	}
	suite.mustSave(&user)

	var loadUser OfficeUser
	err := suite.db.Eager().Find(&loadUser, user.ID)
	suite.Nil(err, "loading user")
	suite.Equal(user.ID, loadUser.ID)
	suite.Equal(office.ID, loadUser.TransportationOffice.ID)
}
