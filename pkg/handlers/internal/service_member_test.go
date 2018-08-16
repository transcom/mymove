package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *utils.HandlerSuite) TestShowServiceMemberHandler() {
	// Given: A servicemember and a user
	user := testdatagen.MakeDefaultUser(suite.db)

	newServiceMember := models.ServiceMember{
		UserID: user.ID,
	}
	suite.mustSave(&newServiceMember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, newServiceMember)

	params := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: str * utils.Fmt.UUID(newServiceMember.ID.String()),
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

func (suite *utils.HandlerSuite) TestShowServiceMemberWrongUser() {
	// Given: Servicemember trying to load another
	notLoggedInUser := testdatagen.MakeDefaultServiceMember(suite.db)
	loggedInUser := testdatagen.MakeDefaultServiceMember(suite.db)

	req := httptest.NewRequest("GET", fmt.Sprintf("/service_members/%s", notLoggedInUser.ID.String()), nil)
	req = suite.authenticateRequest(req, loggedInUser)

	showServiceMemberParams := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: str * utils.Fmt.UUID(notLoggedInUser.ID.String()),
	}
	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := showHandler.Handle(showServiceMemberParams)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.code)
}

func (suite *utils.HandlerSuite) TestSubmitServiceMemberHandlerAllValues() {
	// Given: A logged-in user
	user := testdatagen.MakeDefaultUser(suite.db)

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		UserID:                 str * utils.Fmt.UUID(user.ID.String()),
		Edipi:                  swag.String("random string bla"),
		FirstName:              swag.String("random string bla"),
		MiddleName:             swag.String("random string bla"),
		LastName:               swag.String("random string bla"),
		Suffix:                 swag.String("random string bla"),
		Telephone:              swag.String("random string bla"),
		SecondaryTelephone:     swag.String("random string bla"),
		PersonalEmail:          swag.String("wml@example.com"),
		PhoneIsPreferred:       swag.Bool(false),
		TextMessageIsPreferred: swag.Bool(false),
		EmailIsPreferred:       swag.Bool(true),
		ResidentialAddress:     fakeAddressPayload(),
		BackupMailingAddress:   fakeAddressPayload(),
		SocialSecurityNumber:   (*str * utils.Fmt.SSN)(swag.String("123-45-6789")),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&CookieUpdateResponder{}, response)
	unwrappedResponse := response.(*CookieUpdateResponder).responder
	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)

	// Then: we expect a servicemember to have been created for the user
	query := suite.db.Where(*fmt.Sprintf("user_id='%v'", user.ID))
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.Assertions.Len(serviceMembers, 1)
}

func (suite *utils.HandlerSuite) TestSubmitServiceMemberSSN() {
	// Given: A logged-in user
	user := testdatagen.MakeDefaultUser(suite.db)
	session := &auth.Session{
		UserID:          user.ID,
		ApplicationName: auth.MyApp,
	}

	// When: a new ServiceMember is posted
	ssn := "123-45-6789"
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		SocialSecurityNumber: (*str * utils.Fmt.SSN)(swag.String(ssn)),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.authenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&CookieUpdateResponder{}, response)
	unwrappedResponse := response.(*CookieUpdateResponder).responder
	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)
	okResponse := unwrappedResponse.(*servicememberop.CreateServiceMemberCreated)

	suite.Assertions.True(*okResponse.Payload.HasSocialSecurityNumber)

	// Then: we expect a ServiceMember to have been created for the user
	query := suite.db.Where(*fmt.Sprintf("user_id='%v'", user.ID))
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.Assertions.Len(serviceMembers, 1)

	serviceMemberID, _ := uuid.FromString(okResponse.Payload.ID.String())

	session.ServiceMemberID = serviceMemberID
	serviceMember, err := models.FetchServiceMemberForUser(suite.db, session, serviceMemberID)
	suite.Assertions.NoError(err)

	suite.Assertions.True(serviceMember.SocialSecurityNumber.Matches(ssn))
}

func (suite *utils.HandlerSuite) TestPatchServiceMemberHandler() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.db)

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
	ssn := utils.FmtSSN("555-55-5555")
	resAddress := fakeAddressPayload()
	backupAddress := fakeAddressPayload()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &newEdipi,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &affiliation,
		EmailIsPreferred:     swag.Bool(true),
		FirstName:            swag.String("Firstname"),
		LastName:             swag.String("Lastname"),
		MiddleName:           swag.String("Middlename"),
		PersonalEmail:        swag.String("name@domain.com"),
		PhoneIsPreferred:     swag.Bool(true),
		Rank:                 &rank,
		TextMessageIsPreferred: swag.Bool(true),
		SecondaryTelephone:     swag.String("555555555"),
		SocialSecurityNumber:   ssn,
		Suffix:                 swag.String("Sr."),
		Telephone:              swag.String("555555555"),
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           str * utils.Fmt.UUID(newServiceMember.ID.String()),
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

func (suite *utils.HandlerSuite) TestPatchServiceMemberHandlerWrongUser() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.db)
	user2 := testdatagen.MakeDefaultUser(suite.db)

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
	req = suite.authenticateUserRequest(req, user2)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           str * utils.Fmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.code)
}

func (suite *utils.HandlerSuite) TestPatchServiceMemberHandlerNoServiceMember() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.db)

	servicememberUUID := uuid.Must(uuid.NewV4())

	var newEdipi = "9999999999"

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateUserRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           str * utils.Fmt.UUID(servicememberUUID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.code)
}

func (suite *utils.HandlerSuite) TestPatchServiceMemberHandlerNoChange() {
	// Given: a logged in user with a servicemember
	user := testdatagen.MakeDefaultUser(suite.db)

	var origEdipi = "4444444444"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.mustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.authenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           str * utils.Fmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
}

func (suite *utils.HandlerSuite) TestShowServiceMemberOrders() {
	order1 := testdatagen.MakeDefaultOrder(suite.db)
	order2Assertions := testdatagen.Assertions{
		Order: models.Order{
			ServiceMember:   order1.ServiceMember,
			ServiceMemberID: order1.ServiceMemberID,
		},
	}
	order2 := testdatagen.MakeOrder(suite.db, order2Assertions)

	req := httptest.NewRequest("GET", "/service_members/some_id/current_orders", nil)
	req = suite.authenticateRequest(req, order1.ServiceMember)

	params := servicememberop.ShowServiceMemberOrdersParams{
		HTTPRequest:     req,
		ServiceMemberID: str * utils.Fmt.UUID(order1.ServiceMemberID.String()),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	handler := ShowServiceMemberOrdersHandler(context)

	response := handler.Handle(params)

	suite.IsType(&servicememberop.ShowServiceMemberOrdersOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOrdersOK)
	responsePayload := okResponse.Payload

	// Should return the most recently created order
	suite.Equal(order2.ID.String(), responsePayload.ID.String())
}
