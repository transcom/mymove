package models_test

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicServiceMemberInstantiation() {
	servicemember := &ServiceMember{}

	expErrors := map[string][]string{
		"user_id": {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(servicemember, expErrors)
}

func (suite *ModelSuite) TestIsProfileCompleteWithIncompleteSM() {
	t := suite.T()
	// Given: a user and a service member
	user1 := User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	// And: a service member is incompletely initialized with almost all required values
	edipi := "12345567890"
	affiliation := internalmessages.AffiliationARMY
	rank := internalmessages.ServiceMemberRankE5
	firstName := "bob"
	lastName := "sally"
	telephone := "510 555-5555"
	email := "bobsally@gmail.com"
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
		Affiliation:        &affiliation,
		Rank:               &rank,
		FirstName:          &firstName,
		LastName:           &lastName,
		Telephone:          &telephone,
		PersonalEmail:      &email,
		ResidentialAddress: &fakeAddress,
	}

	// Then: IsProfileComplete should return false
	if servicemember.IsProfileComplete() != false {
		t.Error("Expected profile to be incomplete.")
	}
	// When: all required fields are set
	emailPreferred := true
	servicemember.EmailIsPreferred = &emailPreferred

	// Then: IsProfileComplete should return true
	if servicemember.IsProfileComplete() != true {
		t.Error("Expected profile to be complete.")
	}
}
