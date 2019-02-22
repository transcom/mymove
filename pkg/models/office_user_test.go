package models_test

import (
	"github.com/gofrs/uuid"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_OfficeUserInstantiation() {
	user := &OfficeUser{}
	expErrors := map[string][]string{
		"first_name":               {"FirstName can not be blank."},
		"last_name":                {"LastName can not be blank."},
		"telephone":                {"Telephone can not be blank."},
		"email":                    {"Email can not be blank."},
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
	suite.MustSave(&sally)
	office := CreateTestShippingOffice(suite)

	user := OfficeUser{
		LastName:               "Tester",
		FirstName:              "Sally",
		Email:                  "sally.work@government.gov",
		Telephone:              "(907) 555-1212",
		UserID:                 &sally.ID,
		User:                   sally,
		TransportationOfficeID: office.ID,
	}
	suite.MustSave(&user)

	var loadUser OfficeUser
	err := suite.DB().Eager().Find(&loadUser, user.ID)
	suite.Nil(err, "loading user")
	suite.Equal(user.ID, loadUser.ID)
	suite.Equal(office.ID, loadUser.TransportationOffice.ID)
}

func (suite *ModelSuite) TestFetchOfficeUserByEmail() {
	user, err := FetchOfficeUserByEmail(suite.DB(), "not_here@example.com")
	suite.Equal(err, ErrFetchNotFound)
	suite.Nil(user)

	const email = "sally.work@government.gov"
	office := CreateTestShippingOffice(suite)
	newUser := OfficeUser{
		LastName:               "Tester",
		FirstName:              "Sally",
		Email:                  email,
		Telephone:              "(907) 555-1212",
		TransportationOfficeID: office.ID,
	}
	suite.MustSave(&newUser)

	user, err = FetchOfficeUserByEmail(suite.DB(), email)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(newUser.ID, user.ID)
}

func (suite *ModelSuite) TestFetchOfficeUserByID() {
	fakeUUID, _ := uuid.FromString("99999999-8888-7777-8b57-e39519f42dc1")

	user, err := FetchOfficeUserByID(suite.DB(), fakeUUID)
	suite.NotNil(err)

	office := CreateTestShippingOffice(suite)
	newUser := OfficeUser{
		LastName:               "Tester",
		FirstName:              "Sally",
		Email:                  "test@test.com",
		Telephone:              "(907) 555-1212",
		TransportationOfficeID: office.ID,
	}
	suite.MustSave(&newUser)

	user, err = FetchOfficeUserByID(suite.DB(), newUser.ID)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(newUser.ID, user.ID)
}
