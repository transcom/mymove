package models_test

import (
	"github.com/gofrs/uuid"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_TspUserInstantiation() {
	user := &TspUser{}
	expErrors := map[string][]string{
		"first_name":                         {"FirstName can not be blank."},
		"last_name":                          {"LastName can not be blank."},
		"telephone":                          {"Telephone can not be blank."},
		"email":                              {"Email can not be blank."},
		"transportation_service_provider_id": {"TransportationServiceProviderID can not be blank."},
	}
	suite.verifyValidationErrors(user, expErrors)
}

func (suite *ModelSuite) Test_BasicTspUser() {
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"
	sally := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}
	suite.MustSave(&sally)
	tsp := CreateTestTsp(suite)

	user := TspUser{
		LastName:                        "Tester",
		FirstName:                       "Sally",
		Email:                           "sally.work@government.gov",
		Telephone:                       "(907) 555-1212",
		UserID:                          &sally.ID,
		User:                            sally,
		TransportationServiceProviderID: tsp.ID,
	}
	suite.MustSave(&user)

	var loadUser TspUser
	err := suite.DB().Eager().Find(&loadUser, user.ID)
	suite.Nil(err, "loading user")
	suite.Equal(user.ID, loadUser.ID)
	suite.Equal(tsp.ID, loadUser.TransportationServiceProvider.ID)
}

func (suite *ModelSuite) TestFetchTspUserByEmail() {

	user, err := FetchTspUserByEmail(suite.DB(), "not_here@example.com")
	suite.Equal(err, ErrFetchNotFound)
	suite.Nil(user)

	const email = "sally.work@government.gov"
	tsp := CreateTestTsp(suite)
	newUser := TspUser{
		LastName:                        "Tester",
		FirstName:                       "Sally",
		Email:                           email,
		Telephone:                       "(907) 555-1212",
		TransportationServiceProviderID: tsp.ID,
	}
	suite.MustSave(&newUser)

	user, err = FetchTspUserByEmail(suite.DB(), email)
	suite.Nil(err)
	suite.NotNil(user)
	suite.Equal(newUser.ID, user.ID)
}
