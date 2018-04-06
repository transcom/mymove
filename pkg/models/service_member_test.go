package models_test

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicServiceMemberInstantiation() {
	servicemember := &ServiceMember{}

	expErrors := map[string][]string{
		"user_id": {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(servicemember, expErrors)
}

func (suite *ModelSuite) TestGetServiceMemberForUser() {
	t := suite.T()

	// Given: 2 users
	user1 := User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	user2 := User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "someoneelse@example.com",
	}
	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}
	verrs, err = suite.db.ValidateAndCreate(&user2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	// And: a service member is initialized with an address
	edipi := "12345567890"
	firstName := "bob"
	lastName := "sally"
	telephone := "510 555-5555"
	emailPreferred := true
	fakeAddress := Address{
		StreetAddress1: "123 main st.",
		StreetAddress2: swag.String("Apt.1"),
		City:           "Pleasantville",
		State:          "AL",
		PostalCode:     "01234",
		Country:        swag.String("USA"),
	}

	servicemember := ServiceMember{
		UserID:             user1.ID,
		Edipi:              &edipi,
		FirstName:          &firstName,
		LastName:           &lastName,
		Telephone:          &telephone,
		EmailIsPreferred:   &emailPreferred,
		ResidentialAddress: &fakeAddress,
	}

	// Then: Service Member is created
	verrs, err = suite.db.ValidateAndCreate(&servicemember)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	// And: GetServiceMemberForUser for the existing user returns the service member
	servicememberResult, err := GetServiceMemberForUser(suite.db, user1.ID, servicemember.ID)
	if err != nil {
		t.Error("Expected to get servicememberResult back.", err)
	}
	if !servicememberResult.IsValid() {
		t.Error("Expected the servicemember to be valid")
	}
	if servicememberResult.ServiceMember().ID != servicemember.ID {
		t.Error("Expected new servicemember to match servicemember.")
	}

	// And: GetServiceMemberForUser returns an invalid result for nonexistent SM
	servicememberResult, err = GetServiceMemberForUser(suite.db, user1.ID, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Error("Expected to get servicememberResult back.", err)
	}
	if servicememberResult.IsValid() {
		t.Error("Expected the servicememberResult to be invalid")
	}
	if servicememberResult.ErrorCode() != FetchErrorNotFound {
		t.Error("Should have gotten a not found error")
	}

	// And: GetServiceMemberForUser returns a Forbidden Error for unassociated user
	servicememberResult, err = GetServiceMemberForUser(suite.db, user2.ID, servicemember.ID)
	if err != nil {
		t.Error("Expected to get servicememberResult back.", err)
	}
	if servicememberResult.IsValid() {
		t.Error("Expected the servicememberResult to be invalid")
	}
	if servicememberResult.ErrorCode() != FetchErrorForbidden {
		t.Error("Should have gotten a forbidden error")
	}

}
