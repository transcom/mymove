package models_test

import (
	"strings"

	"github.com/gofrs/uuid"

	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_OfficeUserInstantiation() {
	user := &m.OfficeUser{}
	expErrors := map[string][]string{
		"first_name":               {"FirstName can not be blank."},
		"last_name":                {"LastName can not be blank."},
		"telephone":                {"Telephone can not be blank."},
		"email":                    {"Email does not match the email format."},
		"transportation_office_id": {"TransportationOfficeID can not be blank."},
	}
	suite.verifyValidationErrors(user, expErrors, nil)
}

func (suite *ModelSuite) Test_BasicOfficeUser() {
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"
	sally := m.User{
		OktaID:    fakeUUID.String(),
		OktaEmail: userEmail,
	}
	suite.MustSave(&sally)
	office := CreateTestShippingOffice(suite)

	user := m.OfficeUser{
		LastName:               "Tester",
		FirstName:              "Sally",
		Email:                  "sally.work@government.gov",
		Telephone:              "(907) 555-1212",
		UserID:                 &sally.ID,
		User:                   sally,
		TransportationOfficeID: office.ID,
	}
	suite.MustSave(&user)

	var loadUser m.OfficeUser
	err := suite.DB().Eager().Find(&loadUser, user.ID)
	suite.Nil(err, "loading user")
	suite.Equal(user.ID, loadUser.ID)
	suite.Equal(office.ID, loadUser.TransportationOffice.ID)
}

func (suite *ModelSuite) TestFetchOfficeUserByEmail() {
	user, err := m.FetchOfficeUserByEmail(suite.DB(), "not_here@example.com")
	suite.Equal(err, m.ErrFetchNotFound)
	suite.Nil(user)

	const email = "sally.work@government.gov"
	office := CreateTestShippingOffice(suite)
	newUser := m.OfficeUser{
		LastName:               "Tester",
		FirstName:              "Sally",
		Email:                  email,
		Telephone:              "(907) 555-1212",
		TransportationOfficeID: office.ID,
	}
	suite.MustSave(&newUser)

	user, err = m.FetchOfficeUserByEmail(suite.DB(), email)
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(newUser.ID, user.ID)
}

func (suite *ModelSuite) TestFetchOfficeUserByEmailCaseSensitivity() {
	fakeUUID, _ := uuid.FromString("f390a584-3974-47b9-9ab2-05383304d696")
	userEmail := "Chris@government.gov"

	chris := m.User{
		OktaID:    fakeUUID.String(),
		OktaEmail: userEmail,
	}
	suite.MustSave(&chris)
	office := CreateTestShippingOffice(suite)

	officeUser := m.OfficeUser{
		LastName:               "Tester",
		FirstName:              "Chris",
		Email:                  userEmail,
		Telephone:              "(908) 555-1313",
		UserID:                 &chris.ID,
		User:                   chris,
		TransportationOfficeID: office.ID,
	}
	suite.MustSave(&officeUser)

	user, err := m.FetchOfficeUserByEmail(suite.DB(), strings.ToLower(userEmail))
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(user.Email, userEmail)
}
