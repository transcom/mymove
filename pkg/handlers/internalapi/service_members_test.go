//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package internalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowServiceMemberHandler() {
	// Given: A servicemember and a user
	user := testdatagen.MakeDefaultUser(suite.DB())

	newServiceMember := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID: user.ID,
		},
	})
	suite.MustSave(&newServiceMember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: show ServiceMember is queried
	showHandler := ShowServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&servicememberop.ShowServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOK)

	// And: Returned query to include our added servicemember
	suite.Assertions.Equal(user.ID.String(), okResponse.Payload.UserID.String())
}

func (suite *HandlerSuite) TestShowServiceMemberWrongUser() {
	// Given: Servicemember trying to load another
	notLoggedInUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	loggedInUser := testdatagen.MakeDefaultServiceMember(suite.DB())

	req := httptest.NewRequest("GET", fmt.Sprintf("/service_members/%s", notLoggedInUser.ID.String()), nil)
	req = suite.AuthenticateRequest(req, loggedInUser)

	showServiceMemberParams := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(notLoggedInUser.ID.String()),
	}
	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := showHandler.Handle(showServiceMemberParams)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerNoValues() {
	// Given: A logged-in user
	user := testdatagen.MakeDefaultUser(suite.DB())

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{}

	req := httptest.NewRequest("POST", "/service_members", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}
	handler := CreateServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.CookieUpdateResponder{}, response)

	unwrappedResponse := response.(*handlers.CookieUpdateResponder).Responder
	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)

	// Then: we expect a servicemember to have been created for the user
	query := suite.DB().Where("user_id = ?", user.ID)
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.Assertions.Len(serviceMembers, 1)

	serviceMemberPayload := unwrappedResponse.(*servicememberop.CreateServiceMemberCreated).Payload

	suite.Assertions.NotEqual(*serviceMemberPayload.ID, uuid.Nil)
	suite.Assertions.NotEqual(*serviceMemberPayload.UserID, uuid.Nil)
	suite.Assertions.Equal(*serviceMemberPayload.IsProfileComplete, false)
	suite.Assertions.Equal(len((*serviceMemberPayload).Orders), 0)

	// These shouldn't return any value or Swagger clients will complain during validation
	// because the payloads for these objects are defined to require non-null values for most fields
	// which can't be handled in OpenAPI Spec 2.0. Therefore we don't return them at all.
	suite.Assertions.Equal((*serviceMemberPayload).Rank, (*internalmessages.ServiceMemberRank)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).Affiliation, (*internalmessages.Affiliation)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).CurrentLocation, (*internalmessages.DutyLocationPayload)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).ResidentialAddress, (*internalmessages.Address)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).BackupMailingAddress, (*internalmessages.Address)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).BackupContacts, internalmessages.IndexServiceMemberBackupContactsPayload{})
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerAllValues() {
	// Given: A logged-in user
	user := testdatagen.MakeDefaultUser(suite.DB())

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		UserID:               strfmt.UUID(user.ID.String()),
		Edipi:                swag.String("random string bla"),
		FirstName:            swag.String("random string bla"),
		MiddleName:           swag.String("random string bla"),
		LastName:             swag.String("random string bla"),
		Suffix:               swag.String("random string bla"),
		Telephone:            swag.String("random string bla"),
		SecondaryTelephone:   swag.String("random string bla"),
		PersonalEmail:        swag.String("wml@example.com"),
		PhoneIsPreferred:     swag.Bool(false),
		EmailIsPreferred:     swag.Bool(true),
		ResidentialAddress:   fakeAddressPayload(),
		BackupMailingAddress: fakeAddressPayload(),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.CookieUpdateResponder{}, response)
	unwrappedResponse := response.(*handlers.CookieUpdateResponder).Responder
	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)

	// Then: we expect a servicemember to have been created for the user
	query := suite.DB().Where("user_id = ?", user.ID)
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.Assertions.Len(serviceMembers, 1)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandler() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.DB())

	// TODO: add more fields to change
	var origEdipi = "2342342344"
	var newEdipi = "9999999999"

	origRank := models.ServiceMemberRankE1

	origAffiliation := models.AffiliationAIRFORCE
	newAffiliation := internalmessages.AffiliationARMY

	origFirstName := swag.String("random string bla")
	newFirstName := swag.String("John")

	origMiddleName := swag.String("random string bla")
	newMiddleName := swag.String("")

	origLastName := swag.String("random string bla")
	newLastName := swag.String("Doe")

	origSuffix := swag.String("random string bla")
	newSuffix := swag.String("Mr.")

	origTelephone := swag.String("random string bla")
	newTelephone := swag.String("555-555-5555")

	origSecondaryTelephone := swag.String("random string bla")
	newSecondaryTelephone := swag.String("555-555-5555")

	origPersonalEmail := swag.String("wml@example.com")
	newPersonalEmail := swag.String("example@email.com")

	origPhoneIsPreferred := swag.Bool(false)
	newPhoneIsPreferred := swag.Bool(true)

	origEmailIsPreferred := swag.Bool(true)
	newEmailIsPreferred := swag.Bool(false)

	dutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())

	newServiceMember := models.ServiceMember{
		UserID:             user.ID,
		Edipi:              &origEdipi,
		DutyLocationID:     &dutyLocation.ID,
		DutyLocation:       dutyLocation,
		Rank:               &origRank,
		Affiliation:        &origAffiliation,
		FirstName:          origFirstName,
		MiddleName:         origMiddleName,
		LastName:           origLastName,
		Suffix:             origSuffix,
		Telephone:          origTelephone,
		SecondaryTelephone: origSecondaryTelephone,
		PersonalEmail:      origPersonalEmail,
		PhoneIsPreferred:   origPhoneIsPreferred,
		EmailIsPreferred:   origEmailIsPreferred,
	}
	suite.MustSave(&newServiceMember)

	orderDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
	orderGrade := (string)(models.ServiceMemberRankE5)
	testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMember:        newServiceMember,
			ServiceMemberID:      newServiceMember.ID,
			OriginDutyLocation:   &orderDutyLocation,
			OriginDutyLocationID: &orderDutyLocation.ID,
			Grade:                &orderGrade,
		},
	})

	rank := internalmessages.ServiceMemberRankE1
	resAddress := fakeAddressPayload()
	backupAddress := fakeAddressPayload()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &newEdipi,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &newAffiliation,
		Rank:                 &rank,
		EmailIsPreferred:     newEmailIsPreferred,
		FirstName:            newFirstName,
		LastName:             newLastName,
		MiddleName:           newMiddleName,
		PersonalEmail:        newPersonalEmail,
		PhoneIsPreferred:     newPhoneIsPreferred,
		SecondaryTelephone:   newSecondaryTelephone,
		Suffix:               newSuffix,
		Telephone:            newTelephone,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handlerConfig.SetFileStorer(fakeS3)
	handler := PatchServiceMemberHandler{handlerConfig}
	response := handler.Handle(params)

	suite.IsType(&servicememberop.PatchServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.PatchServiceMemberOK)

	serviceMemberPayload := okResponse.Payload

	suite.Equal(newEdipi, *serviceMemberPayload.Edipi)
	suite.Equal(newAffiliation, *serviceMemberPayload.Affiliation)
	suite.Equal(*newFirstName, *serviceMemberPayload.FirstName)
	suite.Equal(*newMiddleName, *serviceMemberPayload.MiddleName)
	suite.Equal(*newLastName, *serviceMemberPayload.LastName)
	suite.Equal(*newSuffix, *serviceMemberPayload.Suffix)
	suite.Equal(*newTelephone, *serviceMemberPayload.Telephone)
	suite.Equal(*newSecondaryTelephone, *serviceMemberPayload.SecondaryTelephone)
	suite.Equal(*newPersonalEmail, *serviceMemberPayload.PersonalEmail)
	suite.Equal(*newPhoneIsPreferred, *serviceMemberPayload.PhoneIsPreferred)
	suite.Equal(*newEmailIsPreferred, *serviceMemberPayload.EmailIsPreferred)
	suite.Equal(*resAddress.StreetAddress1, *serviceMemberPayload.ResidentialAddress.StreetAddress1)
	suite.Equal(*backupAddress.StreetAddress1, *serviceMemberPayload.BackupMailingAddress.StreetAddress1)
	// Editing SM info DutyLocation and Rank fields should edit Orders OriginDutyLocation and Grade fields
	suite.Equal(*serviceMemberPayload.Orders[0].OriginDutyLocation.Name, newServiceMember.DutyLocation.Name)
	suite.Equal(*serviceMemberPayload.Orders[0].Grade, (string)(rank))
	suite.NotEqual(*serviceMemberPayload.Orders[0].Grade, orderGrade)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerSubmittedMove() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.DB())

	edipi := "2342342344"

	// If there are orders and the move has been submitted, then the
	// affiliation rank, and duty location should not be editable.
	origRank := models.ServiceMemberRankE1
	newRank := internalmessages.ServiceMemberRankE2

	origAffiliation := models.AffiliationAIRFORCE
	newAffiliation := internalmessages.AffiliationARMY

	origDutyLocation := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())
	newDutyLocation := testdatagen.FetchOrMakeDefaultNewOrdersDutyLocation(suite.DB())
	newDutyLocationID := strfmt.UUID(newDutyLocation.ID.String())

	origFirstName := swag.String("random string bla")
	newFirstName := swag.String("John")

	origMiddleName := swag.String("random string bla")
	newMiddleName := swag.String("")

	origLastName := swag.String("random string bla")
	newLastName := swag.String("Doe")

	origSuffix := swag.String("random string bla")
	newSuffix := swag.String("Mr.")

	origTelephone := swag.String("random string bla")
	newTelephone := swag.String("555-555-5555")

	origSecondaryTelephone := swag.String("random string bla")
	newSecondaryTelephone := swag.String("555-555-5555")

	origPersonalEmail := swag.String("wml@example.com")
	newPersonalEmail := swag.String("example@email.com")

	origPhoneIsPreferred := swag.Bool(false)
	newPhoneIsPreferred := swag.Bool(true)

	origEmailIsPreferred := swag.Bool(true)
	newEmailIsPreferred := swag.Bool(false)

	newServiceMember := models.ServiceMember{
		UserID:             user.ID,
		Edipi:              &edipi,
		Rank:               &origRank,
		Affiliation:        &origAffiliation,
		DutyLocationID:     &origDutyLocation.ID,
		DutyLocation:       origDutyLocation,
		FirstName:          origFirstName,
		MiddleName:         origMiddleName,
		LastName:           origLastName,
		Suffix:             origSuffix,
		Telephone:          origTelephone,
		SecondaryTelephone: origSecondaryTelephone,
		PersonalEmail:      origPersonalEmail,
		PhoneIsPreferred:   origPhoneIsPreferred,
		EmailIsPreferred:   origEmailIsPreferred,
	}
	suite.MustSave(&newServiceMember)

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: newServiceMember,
	})

	// The testdatagen sets these values, fails if you try to blank them out via Assertions,
	// and gives defaults if you pass nil, so we have to set this after the creation.
	// This more closely resembles what orders would look like pre and post submission, before
	// a TOO gets to them.
	move.Orders.TAC = nil
	move.Orders.DepartmentIndicator = nil
	move.Orders.OrdersNumber = nil
	move.Orders.OrdersTypeDetail = nil

	suite.MustSave(&move.Orders)
	moveRouter := moverouter.NewMoveRouter()
	moveRouter.Submit(suite.AppContextForTest(), &move)
	suite.MustSave(&move)

	resAddress := fakeAddressPayload()
	backupAddress := fakeAddressPayload()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &edipi,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &newAffiliation,
		EmailIsPreferred:     newEmailIsPreferred,
		FirstName:            newFirstName,
		LastName:             newLastName,
		MiddleName:           newMiddleName,
		PersonalEmail:        newPersonalEmail,
		PhoneIsPreferred:     newPhoneIsPreferred,
		Rank:                 &newRank,
		SecondaryTelephone:   newSecondaryTelephone,
		Suffix:               newSuffix,
		Telephone:            newTelephone,
		CurrentLocationID:    &newDutyLocationID,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := PatchServiceMemberHandler{handlerConfig}
	response := handler.Handle(params)

	suite.IsType(&servicememberop.PatchServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.PatchServiceMemberOK)

	serviceMemberPayload := okResponse.Payload

	// These fields should not change (they should still be the original
	// values) after the move has been submitted.
	suite.Equal(origAffiliation, models.ServiceMemberAffiliation(*serviceMemberPayload.Affiliation))
	suite.Equal(origRank, models.ServiceMemberRank(*serviceMemberPayload.Rank))
	suite.Equal(origDutyLocation.ID.String(), string(*serviceMemberPayload.CurrentLocation.ID))

	// These fields should change even if the move is submitted.
	suite.Equal(*newFirstName, *serviceMemberPayload.FirstName)
	suite.Equal(*newMiddleName, *serviceMemberPayload.MiddleName)
	suite.Equal(*newLastName, *serviceMemberPayload.LastName)
	suite.Equal(*newSuffix, *serviceMemberPayload.Suffix)
	suite.Equal(*newTelephone, *serviceMemberPayload.Telephone)
	suite.Equal(*newSecondaryTelephone, *serviceMemberPayload.SecondaryTelephone)
	suite.Equal(*newPersonalEmail, *serviceMemberPayload.PersonalEmail)
	suite.Equal(*newPhoneIsPreferred, *serviceMemberPayload.PhoneIsPreferred)
	suite.Equal(*newEmailIsPreferred, *serviceMemberPayload.EmailIsPreferred)

	suite.Equal(*resAddress.StreetAddress1, *serviceMemberPayload.ResidentialAddress.StreetAddress1)
	suite.Equal(*backupAddress.StreetAddress1, *serviceMemberPayload.BackupMailingAddress.StreetAddress1)

	// Then: we expect addresses to have been created
	addresses := []models.Address{}
	suite.DB().All(&addresses)
	suite.Equal(8, len(addresses))
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerWrongUser() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.DB())
	user2 := testdatagen.MakeDefaultUser(suite.DB())

	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.MustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateUserRequest(req, user2)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoServiceMember() {
	// Given: a logged in user
	user := testdatagen.MakeDefaultUser(suite.DB())

	servicememberUUID := uuid.Must(uuid.NewV4())

	var newEdipi = "9999999999"

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(servicememberUUID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoChange() {
	// Given: a logged in user with a servicemember
	user := testdatagen.MakeDefaultUser(suite.DB())

	var origEdipi = "4444444444"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.MustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler{handlers.NewHandlerConfig(suite.DB(), suite.Logger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
}

func (suite *HandlerSuite) TestShowServiceMemberOrders() {
	order1 := testdatagen.MakeDefaultOrder(suite.DB())
	order2Assertions := testdatagen.Assertions{
		Order: models.Order{
			ServiceMember:   order1.ServiceMember,
			ServiceMemberID: order1.ServiceMemberID,
		},
	}
	order2 := testdatagen.MakeOrder(suite.DB(), order2Assertions)

	req := httptest.NewRequest("GET", "/service_members/some_id/current_orders", nil)
	req = suite.AuthenticateRequest(req, order1.ServiceMember)

	params := servicememberop.ShowServiceMemberOrdersParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(order1.ServiceMemberID.String()),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handlerConfig.SetFileStorer(fakeS3)
	handler := ShowServiceMemberOrdersHandler{handlerConfig}

	response := handler.Handle(params)

	suite.IsType(&servicememberop.ShowServiceMemberOrdersOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOrdersOK)
	responsePayload := okResponse.Payload

	// Should return the most recently created order
	suite.Equal(order2.ID.String(), responsePayload.ID.String())
}
