package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestShowServiceMemberHandler() {
	// Given: A servicemember and a user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	newServiceMember := models.ServiceMember{
		UserID: user.ID,
	}
	suite.mustSave(&newServiceMember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user)

	params := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: show ServiceMember is queried
	showHandler := ShowServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&servicememberop.ShowServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOK)

	// And: Returned query to include our added servicemember
	suite.Assertions.Equal(user.ID.String(), okResponse.Payload.UserID.String())
}

func (suite *HandlerSuite) TestShowServiceMemberWrongUser() {
	// Given: A servicemember with a not-logged-in user and a separate logged-in user
	notLoggedInUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&notLoggedInUser)

	loggedInUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&loggedInUser)

	// When: A servicemember is created for not-logged-in-user
	newServiceMember := models.ServiceMember{
		UserID: notLoggedInUser.ID,
	}
	suite.mustSave(&newServiceMember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, loggedInUser)

	showServiceMemberParams := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := showHandler.Handle(showServiceMemberParams)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.code)
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerAllValues() {
	// Given: A logged-in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		UserID:                 strfmt.UUID(user.ID.String()),
		Edipi:                  swag.String("random string bla"),
		FirstName:              swag.String("random string bla"),
		MiddleName:             swag.String("random string bla"),
		LastName:               swag.String("random string bla"),
		Suffix:                 swag.String("random string bla"),
		Telephone:              swag.String("random string bla"),
		SecondaryTelephone:     swag.String("random string bla"),
		PersonalEmail:          fmtEmail("wml@example.com"),
		PhoneIsPreferred:       swag.Bool(false),
		TextMessageIsPreferred: swag.Bool(false),
		EmailIsPreferred:       swag.Bool(true),
		ResidentialAddress:     fakeAddress(),
		BackupMailingAddress:   fakeAddress(),
		SocialSecurityNumber:   (*strfmt.SSN)(swag.String("123-45-6789")),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, response)

	// Then: we expect a servicemember to have been created for the user
	query := suite.db.Where(fmt.Sprintf("user_id='%v'", user.ID))
	servicemembers := []models.ServiceMember{}
	query.All(&servicemembers)

	suite.Assertions.Len(servicemembers, 1)
}

func (suite *HandlerSuite) TestSubmitServiceMemberSSN() {
	// Given: A logged-in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// When: a new ServiceMember is posted
	ssn := "123-45-6789"
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		SocialSecurityNumber: (*strfmt.SSN)(swag.String(ssn)),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, response)
	okResponse := response.(*servicememberop.CreateServiceMemberCreated)

	suite.Assertions.True(*okResponse.Payload.HasSocialSecurityNumber)

	// Then: we expect a servicemember to have been created for the user
	query := suite.db.Where(fmt.Sprintf("user_id='%v'", user.ID))
	servicemembers := []models.ServiceMember{}
	query.All(&servicemembers)

	suite.Assertions.Len(servicemembers, 1)

	serviceMemberID, _ := uuid.FromString(okResponse.Payload.ID.String())
	serviceMember, err := models.FetchServiceMember(suite.db, user, serviceMemberID)
	suite.Assertions.NoError(err)

	suite.Assertions.True(serviceMember.SocialSecurityNumber.Matches(ssn))
}

func (suite *HandlerSuite) TestPatchServiceMemberHandler() {
	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// TODO: add more fields to change
	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.mustSave(&newServiceMember)

	affiliation := internalmessages.AffiliationARMY
	rank := internalmessages.ServiceMemberRankE1
	ssn := fmtSSN("555-55-5555")
	resAddress := fakeAddress()
	backupAddress := fakeAddress()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &newEdipi,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &affiliation,
		EmailIsPreferred:     swag.Bool(true),
		FirstName:            swag.String("Firstname"),
		LastName:             swag.String("Lastname"),
		MiddleName:           swag.String("Middlename"),
		PersonalEmail:        fmtEmail("name@domain.com"),
		PhoneIsPreferred:     swag.Bool(true),
		Rank:                 &rank,
		TextMessageIsPreferred: swag.Bool(true),
		SecondaryTelephone:     swag.String("555555555"),
		SocialSecurityNumber:   ssn,
		Suffix:                 swag.String("Sr."),
		Telephone:              swag.String("555555555"),
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.PatchServiceMemberOK)

	serviceMemberPayload := okResponse.Payload

	suite.Assertions.Equal(*serviceMemberPayload.Edipi, newEdipi)
	suite.Assertions.Equal(*serviceMemberPayload.Affiliation, affiliation)
	suite.Assertions.Equal(*serviceMemberPayload.HasSocialSecurityNumber, true)
	suite.Assertions.Equal(*serviceMemberPayload.TextMessageIsPreferred, true)
	suite.Assertions.Equal(*serviceMemberPayload.ResidentialAddress.StreetAddress1, *resAddress.StreetAddress1)
	suite.Assertions.Equal(*serviceMemberPayload.BackupMailingAddress.StreetAddress1, *backupAddress.StreetAddress1)

	// Then: we expect addresses to have been created
	addresses := []models.Address{}
	suite.db.All(&addresses)
	suite.Assertions.Len(addresses, 2)

	// Then: we expect a SSN to have been created
	ssns := models.SocialSecurityNumbers{}
	suite.db.All(&ssns)
	suite.Assertions.Len(ssns, 1)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerWrongUser() {
	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&user2)

	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.mustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user2)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoServiceMember() {
	// Given: a logged in user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	servicememberUUID := uuid.Must(uuid.NewV4())

	var newEdipi = "9999999999"

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(servicememberUUID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoChange() {
	// Given: a logged in user with a servicemember
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	var origEdipi = "4444444444"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.mustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
}
