package models_test

import (
	"github.com/gobuffalo/uuid"

	"github.com/markbates/goth"
	. "github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func (suite *ModelSuite) TestUserCreation() {
	t := suite.T()

	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"

	newUser := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}

	if err := suite.db.Create(&newUser); err != nil {
		t.Fatal("Didn't create user in db.")
	}

	if newUser.ID == uuid.Nil {
		t.Error("Didn't get an id back for user.")
	}

	if (newUser.LoginGovEmail != userEmail) &&
		(newUser.LoginGovUUID != fakeUUID) {
		t.Error("Required values didn't get set.")
	}
}

func (suite *ModelSuite) TestUserCreationWithoutValues() {
	newUser := &User{}

	expErrors := map[string][]string{
		"login_gov_email": {"LoginGovEmail can not be blank."},
		"login_gov_uuid":  {"LoginGovUUID can not be blank."},
	}

	suite.verifyValidationErrors(newUser, expErrors)
}

func (suite *ModelSuite) TestUserCreationDuplicateUUID() {
	t := suite.T()

	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	userEmail := "sally@government.gov"

	newUser := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}

	sameUser := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}

	suite.db.Create(&newUser)
	err := suite.db.Create(&sameUser)

	if err.Error() != `pq: duplicate key value violates unique constraint "constraint_name"` {
		t.Fatal("Db should have errored on unique constraint for UUID")
	}
}

func (suite *ModelSuite) TestGetOrCreateUser() {
	t := suite.T()

	// When: login gov UUID is passed to create user func
	gothUser := goth.User{Email: "sally@government.gov", UserID: "39b28c92-0506-4bef-8b57-e39519f42dc2"}
	loginGovUUID, _ := uuid.FromString(gothUser.UserID)

	// And: user does not yet exist in the db
	newUser, err := GetOrCreateUser(suite.db, gothUser)
	if err != nil {
		t.Error("error querying or creating user.", err)
	}

	// Then: expect fields to be set on returned user
	if newUser.LoginGovEmail != gothUser.Email {
		t.Error("expected email to be set")
	}
	if newUser.LoginGovUUID != loginGovUUID {
		t.Error("expected uuid to be set")
	}

	// When: The same UUID is passed in func
	sameUser, err := GetOrCreateUser(suite.db, gothUser)
	if err != nil {
		t.Error("error querying or creating user.")
	}

	// Then: expect the existing user to be returned
	if sameUser.LoginGovEmail != newUser.LoginGovEmail {
		t.Error("expected existing user to have been returned")
	}

	// And: no new user to have been created
	query := suite.db.Where("login_gov_uuid = $1", loginGovUUID)
	var users []User
	queryErr := query.All(&users)
	if queryErr != nil {
		t.Error("DB Query Error", zap.Error(err))
	}
	if len(users) > 1 {
		t.Error("1 user should have been returned")
	}
}
